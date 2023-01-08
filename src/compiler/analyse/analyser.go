package analyse

import (
	"github.com/kkkunny/Sim/src/compiler/parse"
	"github.com/kkkunny/Sim/src/compiler/utils"
	"github.com/kkkunny/stl/list"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/set"
	"github.com/kkkunny/stl/types"
)

// *********************************************************************************************************************

// AnalyseMain 作为主包进行语义分析
func AnalyseMain(ast *parse.Package) (*ProgramContext, error) {
	ctx := newProgramContext()
	// 包
	pkgCtx := newPackageContext(ctx, ast.Path)
	ctx.importedPackageSet[ast.Path] = pkgCtx
	if err := analyseNoMain(pkgCtx, ast); err != nil {
		return nil, err
	}
	return ctx, nil
}

// 作为辅包进行语义分析
func analyseNoMain(ctx *packageContext, ast *parse.Package) error {
	// 包导入
	rootPath, err := utils.GetRootPath()
	if err != nil {
		return err
	}
	for _, fileAst := range ast.Files {
		for iter := fileAst.Globals.Iterator(); iter.HasValue(); iter.Next() {
			// 获取包路径
			importAst, ok := iter.Value().(*parse.Import)
			if !ok {
				continue
			}
			var pkgPath stlos.Path
			for _, p := range importAst.Packages {
				pkgPath = pkgPath.Join(stlos.Path(p.Source))
			}
			pkgPath = rootPath.Join(pkgPath)
			if !pkgPath.IsExist() {
				return utils.Errorf(importAst.Position(), "unknown package `%s`", pkgPath)
			}
			// 包名
			var pkgName string
			var pkgPos utils.Position
			if importAst.Suffix == nil {
				pkgName = pkgPath.GetBase().String()
				pkgPos = importAst.Packages[len(importAst.Packages)-1].Pos
			} else if importAst.Suffix.IsRight() {
				pkgName = importAst.Suffix.Right().Source
				pkgPos = importAst.Suffix.Right().Pos
			}
			// 导入包
			if pkgCtx, ok := ctx.f.importedPackageSet[pkgPath]; !ok {
				// 从没导入过
				if importAst.Suffix != nil && importAst.Suffix.IsLeft() {
					ctx.f.importedPackageSet[pkgPath] = nil
					pkgCtx = newPackageContext(ctx.f, pkgPath)
					ctx.includes = append(ctx.includes, pkgCtx)
				} else {
					if _, ok := ctx.externs[pkgName]; ok {
						return utils.Errorf(pkgPos, "duplicate identifier")
					}
					ctx.f.importedPackageSet[pkgPath] = nil
					pkgCtx = newPackageContext(ctx.f, pkgPath)
					ctx.externs[pkgName] = pkgCtx
				}
				// 语法分析
				pkgAst, err := parse.ParsePackage(pkgPath)
				if err != nil {
					return err
				}
				if err = analyseNoMain(pkgCtx, pkgAst); err != nil {
					return err
				}
				ctx.f.importedPackageSet[pkgPath] = pkgCtx
			} else {
				// 以前导入过
				if pkgCtx == nil {
					return utils.Errorf(importAst.Position(), "circular reference package `%s`", pkgPath)
				}
				if importAst.Suffix != nil && importAst.Suffix.IsLeft() {
					ctx.includes = append(ctx.includes, pkgCtx)
				} else {
					if c, ok := ctx.externs[pkgName]; ok && pkgCtx != c {
						return utils.Errorf(pkgPos, "duplicate identifier")
					}
					ctx.externs[pkgName] = pkgCtx
				}
			}
		}
	}
	// 包体
	if err := analysePackage(ctx, ast); err != nil {
		return err
	}
	return nil
}

// 包
func analysePackage(ctx *packageContext, ast *parse.Package) utils.Error {
	// 类型定义
	for _, file := range ast.Files {
		err := analysePackageTypeDef(ctx, file.Globals)
		if err != nil {
			return err
		}
	}
	// 变量声明
	for _, file := range ast.Files {
		err := analysePackageVariableDecl(ctx, file.Globals)
		if err != nil {
			return err
		}
	}
	// 变量定义
	for _, file := range ast.Files {
		err := analysePackageVariableDef(ctx, file.Globals)
		if err != nil {
			return err
		}
	}
	return nil
}

// 包 类型定义
func analysePackageTypeDef(ctx *packageContext, asts *list.SingleLinkedList[parse.Global]) utils.Error {
	var errors []utils.Error
	typedefs := list.NewSingleLinkedList[*parse.TypeDef]()
	// 定义
	for iter := asts.Iterator(); iter.HasValue(); iter.Next() {
		ast, ok := iter.Value().(*parse.TypeDef)
		if !ok {
			continue
		}
		if _, ok := ctx.typedefs[ast.Name.Source]; ok {
			errors = append(errors, utils.Errorf(ast.Name.Pos, "duplicate identifier"))
			continue
		}

		ctx.typedefs[ast.Name.Source] = types.NewPair(ast.Public, NewTypedef(ctx.path, ast.Name.Source, nil))
		typedefs.Add(ast)
	}
	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 1 {
		return utils.NewMultiError(errors...)
	}
	// 解析目标类型
	for iter := typedefs.Iterator(); iter.HasValue(); iter.Next() {
		dst, err := analyseType(ctx, iter.Value().Target)
		if err != nil {
			errors = append(errors, err)
		} else {
			ctx.typedefs[iter.Value().Name.Source].Second.Dst = dst
		}
	}
	// 循环引用检测
	for iter := typedefs.Iterator(); iter.HasValue(); iter.Next() {
		ast := iter.Value()
		if checkTypeCircle(set.NewLinkedHashSet[*Typedef](), ctx.typedefs[ast.Name.Source].Second) {
			errors = append(errors, utils.Errorf(ast.Name.Pos, "circular reference"))
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

// 包 变量声明
func analysePackageVariableDecl(ctx *packageContext, asts *list.SingleLinkedList[parse.Global]) utils.Error {
	var errors []utils.Error
	for iter := asts.Iterator(); iter.HasValue(); iter.Next() {
		var g Global
		var err utils.Error
		switch global := iter.Value().(type) {
		case *parse.ExternFunction:
			g, err = analyseExternFunction(ctx, global)
		case *parse.Function:
			g, err = analyseFunctionDecl(ctx, global)
		case *parse.Method:
			g, err = analyseMethodDecl(ctx, global)
		case *parse.GlobalValue:
			g, err = analyseGlobalVariable(ctx, global)
		default:
			continue
		}
		if err != nil {
			errors = append(errors, err)
		} else {
			ctx.f.Globals = append(ctx.f.Globals, g)
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

// 包 变量定义
func analysePackageVariableDef(ctx *packageContext, asts *list.SingleLinkedList[parse.Global]) utils.Error {
	var errors []utils.Error
	for iter := asts.Iterator(); iter.HasValue(); iter.Next() {
		switch global := iter.Value().(type) {
		case *parse.Function:
			if global.Body == nil {
				continue
			}
			if err := analyseFunctionDef(ctx, ctx.GetValue(global.Name.Source).Second.(*Function), global); err != nil {
				errors = append(errors, err)
			}
		case *parse.Method:
			if err := analyseMethodDef(ctx, global); err != nil {
				errors = append(errors, err)
			}
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
