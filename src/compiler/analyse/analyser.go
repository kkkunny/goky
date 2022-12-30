package analyse

import (
	"github.com/kkkunny/klang/src/compiler/utils"
	"github.com/kkkunny/stl/list"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/set"
	"github.com/kkkunny/stl/types"
)

// *********************************************************************************************************************

// AnalyseMain 作为主包进行语义分析
func AnalyseMain(ast parse.Package) (*ProgramContext, error) {
	ctx := newProgramContext()
	// 包
	pkgCtx := newPackageContext(ctx, ast.PkgPath)
	ctx.importedPackageSet[ast.PkgPath] = pkgCtx
	if err := analyseNoMain(pkgCtx, ast); err != nil {
		return nil, err
	}
	// 模板
	for _, pkg := range ctx.importedPackageSet {
		for _, ft := range pkg.funcTemplates {
			for _, f := range ft.Second.impls {
				ctx.Globals = append(ctx.Globals, f)
			}
		}
	}
	return ctx, nil
}

// 作为辅包进行语义分析
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
	return nil
}

// 包
func analysePackage(ctx *packageContext, ast parse.Package) utils.Error {
	// 类型定义
	err := analysePackageTypeDef(ctx, ast.Globals)
	if err != nil {
		return err
	}
	// 变量声明
	err = analysePackageVariableDecl(ctx, ast.Globals)
	if err != nil {
		return err
	}
	// 变量定义
	err = analysePackageVariableDef(ctx, ast.Globals)
	return err
}

// 包 类型定义
func analysePackageTypeDef(ctx *packageContext, asts []parse.Global) utils.Error {
	var errors []utils.Error
	typedefs := list.NewSingleLinkedList[*parse.Typedef]()
	// 定义
	for _, ag := range asts {
		if ag.GlobalNoAttr == nil || ag.GlobalNoAttr.TypeDef == nil {
			continue
		}

		ast := ag.GlobalNoAttr.TypeDef
		if _, ok := ctx.typedefs[ast.Name.Value]; ok {
			errors = append(errors, utils.Errorf(ast.Name.Position, "duplicate identifier"))
			continue
		}

		if len(ast.Templates) == 0 {
			ctx.typedefs[ast.Name.Value] = types.NewPair(ast.Public != nil, NewTypedef(ctx.path, ast.Name.Value, nil))
			typedefs.Add(ast)
		} else {
			ctx.typedefTemplates[ast.Name.Value] = types.NewPair(ast.Public != nil, &typedefTemplate{
				pkg:   ctx.path,
				ast:   ast,
				impls: make(map[string]*Typedef),
			})
		}
	}
	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 1 {
		return utils.NewMultiError(errors...)
	}
	// 解析目标类型
	for iter := typedefs.Iterator(); iter.HasValue(); iter.Next() {
		if len(iter.Value().Templates) != 0 {
			continue
		}
		dst, err := analyseType(ctx, &iter.Value().Dst)
		if err != nil {
			errors = append(errors, err)
		} else {
			ctx.typedefs[iter.Value().Name.Value].Second.Dst = dst
		}
	}
	// 循环引用检测
	for iter := typedefs.Iterator(); iter.HasValue(); iter.Next() {
		ast := iter.Value()
		if checkTypeCircle(set.NewLinkedHashSet[*Typedef](), ctx.typedefs[ast.Name.Value].Second) {
			errors = append(errors, utils.Errorf(ast.Name.Position, "circular reference"))
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
func analysePackageVariableDecl(ctx *packageContext, asts []parse.Global) utils.Error {
	var errors []utils.Error
	for _, ag := range asts {
		if ag.GlobalWithAttr == nil {
			continue
		}
		switch {
		case ag.GlobalWithAttr.Global.Function != nil:
			fast := ag.GlobalWithAttr.Global.Function
			if fast.Tail.Function != nil {
				if len(fast.Tail.Function.Templates) == 0 {
					f, err := analyseFunctionDecl(ctx, ag.GlobalWithAttr.Attr, *fast)
					if err != nil {
						errors = append(errors, err)
					} else {
						ctx.f.Globals = append(ctx.f.Globals, f)
					}
				} else {
					if err := analyseFunctionTemplateDecl(ctx, ag.GlobalWithAttr.Attr, fast); err != nil {
						errors = append(errors, err)
					}
				}
			} else {
				f, err := analyseMethodDecl(ctx, ag.GlobalWithAttr.Attr, *fast)
				if err != nil {
					errors = append(errors, err)
				} else {
					ctx.f.Globals = append(ctx.f.Globals, f)
				}
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
	if len(errors) == 0 {
		return nil
	} else if len(errors) == 1 {
		return errors[0]
	} else {
		return utils.NewMultiError(errors...)
	}
}

// 包 变量定义
func analysePackageVariableDef(ctx *packageContext, asts []parse.Global) utils.Error {
	var errors []utils.Error
	for _, ag := range asts {
		if ag.GlobalWithAttr == nil {
			continue
		}
		switch {
		case ag.GlobalWithAttr.Global.Function != nil:
			fast := ag.GlobalWithAttr.Global.Function
			if fast.Tail.Function != nil {
				if fast.Tail.Function.Body == nil || len(fast.Tail.Function.Templates) != 0 {
					continue
				}
				if err := analyseFunctionDef(ctx, ctx.GetValue(fast.Tail.Function.Name.Value).Second.(*Function), *fast.Tail.Function); err != nil {
					errors = append(errors, err)
				}
			} else {
				if err := analyseMethodDef(ctx, *fast.Tail.Method); err != nil {
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
