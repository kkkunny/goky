package analyse

import (
	"github.com/kkkunny/klang/src/compiler/internal/parse"
	"github.com/kkkunny/klang/src/compiler/internal/utils"
	stlos "github.com/kkkunny/stl/os"
)

// *********************************************************************************************************************

// AnalyseMain 作为主包进行语义分析
func AnalyseMain(ast parse.Package) (*ProgramContext, error) {
	ctx := newProgramContext()
	pkgCtx := newPackageContext(ctx, ast.PkgPath)
	ctx.importedPackageSet[ast.PkgPath] = pkgCtx
	if err := analyseNoMain(pkgCtx, ast); err != nil {
		return nil, err
	}
	if pkgCtx.main != nil {
		pkgCtx.main.Main = true
	}
	return ctx, nil
}

// analyseNoMain 作为辅包进行语义分析
func analyseNoMain(ctx *packageContext, ast parse.Package) error {
	// 包导入
	if len(ast.Globals) > 0 {
		rootPath, err := utils.GetRootPath()
		if err != nil {
			return err
		}
		for _, ag := range ast.Globals {
			// 获取包路径
			if ag.GlobalNoAttr == nil || ag.GlobalNoAttr.Import == nil {
				continue
			}
			var pkgPath stlos.Path
			for _, p := range ag.GlobalNoAttr.Import.Packages {
				pkgPath = pkgPath.Join(stlos.Path(p.Value))
			}
			pkgPath = rootPath.Join(pkgPath)
			if !pkgPath.IsExist() {
				return utils.Errorf(ag.GlobalNoAttr.Import.Position, "unknown package `%s`", pkgPath)
			}
			// 包名
			var pkgName string
			var pkgPos utils.Position
			if ag.GlobalNoAttr.Import.Alias == nil {
				pkgName = pkgPath.GetBase().String()
				pkgPos = ag.GlobalNoAttr.Import.Packages[len(ag.GlobalNoAttr.Import.Packages)-1].Position
			} else {
				pkgName = ag.GlobalNoAttr.Import.Alias.Value
				pkgPos = ag.GlobalNoAttr.Import.Alias.Position
			}
			if pkgCtx, ok := ctx.f.importedPackageSet[pkgPath]; !ok {
				if _, ok := ctx.externs[pkgName]; ok {
					return utils.Errorf(pkgPos, "duplicate identifier")
				}
				ctx.f.importedPackageSet[pkgPath] = nil
				pkgCtx := newPackageContext(ctx.f, pkgPath)
				ctx.externs[pkgName] = pkgCtx
				// 语法分析
				ast, err := parse.ParsePackage(pkgPath)
				if err != nil {
					return err
				}
				if err = analyseNoMain(pkgCtx, *ast); err != nil {
					return err
				}
				ctx.f.importedPackageSet[pkgPath] = pkgCtx
			} else {
				if pkgCtx == nil {
					return utils.Errorf(ag.GlobalNoAttr.Import.Position, "circular reference package `%s`", pkgPath)
				}
				if c, ok := ctx.externs[pkgName]; ok && pkgCtx != c {
					return utils.Errorf(pkgPos, "duplicate identifier")
				}
				ctx.externs[pkgName] = pkgCtx
			}
		}
	}
	// 包体
	if err := analysePackage(ctx, ast); err != nil {
		return err
	}
	// 屏蔽包主函数
	if ctx.main != nil {
		ctx.main.Main = false
	}
	return nil
}

// 包
func analysePackage(ctx *packageContext, ast parse.Package) utils.Error {
	var errors []utils.Error
	// 函数声明
	for _, ag := range ast.Globals {
		if ag.GlobalWithAttr == nil {
			continue
		}
		switch {
		case ag.GlobalWithAttr.Global.Function != nil:
			f, err := analyseFunctionDecl(ctx, ag.GlobalWithAttr.Attr, *ag.GlobalWithAttr.Global.Function)
			if err != nil {
				errors = append(errors, err)
			} else {
				ctx.f.Globals = append(ctx.f.Globals, f)
			}
		case ag.GlobalWithAttr.Global.Variable != nil:
			v, err := analyseGlobalVariable(ctx, ag.GlobalWithAttr.Attr, *ag.GlobalWithAttr.Global.Variable)
			if err != nil {
				errors = append(errors, err)
			} else {
				ctx.f.Globals = append(ctx.f.Globals, v)
			}
		default:
			panic("")
		}
	}

	// 函数定义
	for _, ag := range ast.Globals {
		if ag.GlobalWithAttr == nil {
			continue
		}
		switch {
		case ag.GlobalWithAttr.Global.Function != nil:
			if ag.GlobalWithAttr.Global.Function.Body != nil {
				_, err := analyseFunctionDef(ctx, *ag.GlobalWithAttr.Global.Function)
				if err != nil {
					errors = append(errors, err)
				}
			}
		case ag.GlobalWithAttr.Global.Variable != nil:
		default:
			panic("")
		}
	}

	if len(errors) == 0 {
		return nil
	} else if len(errors) == 1 {
		return errors[0]
	} else {
		return utils.NewMultiError(errors...)
	}
}
