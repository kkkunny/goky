package analyse

import (
	"github.com/kkkunny/klang/src/compiler/parse"
	"github.com/kkkunny/klang/src/compiler/utils"
	stlos "github.com/kkkunny/stl/os"
)

// Global 全局
type Global interface {
	global()
}

// Function 函数
type Function struct {
	// 属性
	ExternName string // 外部名
	NoReturn   bool   // 函数是否不返回
	Exit       bool   // 函数是否会导致程序退出

	Ret    Type
	Params []*Param
	Body   *Block
}

func (self Function) global() {}

func (self Function) stmt() {}

func (self Function) ident() {}

func (self Function) GetType() Type {
	paramTypes := make([]Type, len(self.Params))
	for i, p := range self.Params {
		paramTypes[i] = p.Type
	}
	return NewFuncType(self.Ret, paramTypes...)
}

func (self Function) GetMut() bool {
	return false
}

func (self Function) IsTemporary() bool {
	return true
}

// GlobalVariable 全局变量
type GlobalVariable struct {
	ExternName string

	Type  Type
	Value Expr
}

func (self GlobalVariable) global() {}

func (self GlobalVariable) stmt() {}

func (self GlobalVariable) ident() {}

func (self GlobalVariable) GetType() Type {
	return self.Type
}

func (self GlobalVariable) GetMut() bool {
	return true
}

func (self GlobalVariable) IsTemporary() bool {
	return false
}

// *********************************************************************************************************************

// 函数声明
func analyseFunctionDecl(ctx *packageContext, astAttrs []parse.Attr, ast parse.Function) (*Function, utils.Error) {
	retType, err := analyseType(ctx, ast.Ret)
	if err != nil {
		return nil, err
	}

	params := make([]*Param, len(ast.Params))
	var errors []utils.Error
	for i, p := range ast.Params {
		pt, err := analyseType(ctx, &p.Type)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		params[i] = &Param{Type: pt}
	}
	if len(errors) == 1 {
		return nil, errors[0]
	} else if len(errors) > 1 {
		return nil, utils.NewMultiError(errors...)
	}

	f := &Function{
		Ret:    retType,
		Params: params,
	}

	// 属性
	errors = make([]utils.Error, 0)
	for _, astAttr := range astAttrs {
		switch {
		case astAttr.Extern != nil:
			f.ExternName = astAttr.Extern.Value
		case astAttr.LinkAsm != nil:
			linkPath := stlos.Path(*astAttr.LinkAsm)
			if !linkPath.IsAbsolute() {
				linkPath = ctx.path.Join(linkPath)
			}
			if !linkPath.IsExist() {
				errors = append(errors, utils.Errorf(astAttr.Position, "can not find path `%s`", linkPath))
			}
			ctx.f.Links[linkPath] = struct{}{}
		case astAttr.LinkLib != nil:
			ctx.f.Libs[string(*astAttr.LinkLib)] = struct{}{}
		case astAttr.NoReturn != nil:
			f.NoReturn = true
		case astAttr.Exit != nil:
			f.Exit = true
			f.NoReturn = true
		default:
			panic("")
		}
	}
	if f.ExternName == "" && ast.Body == nil {
		errors = append(errors, utils.Errorf(ast.Name.Position, "missing function body"))
	}
	if len(errors) == 1 {
		return nil, errors[0]
	} else if len(errors) > 1 {
		return nil, utils.NewMultiError(errors...)
	}

	if !ctx.AddValue(ast.Public != nil, ast.Name.Value, f) {
		return nil, utils.Errorf(ast.Name.Position, "duplicate identifier")
	}
	return f, nil
}

// 函数定义
func analyseFunctionDef(ctx *packageContext, ast parse.Function) (*Function, utils.Error) {
	obj := ctx.GetValue(ast.Name.Value).Second
	if obj == nil {
		return nil, utils.Errorf(ast.Name.Position, "unknown identifier")
	}
	f, ok := obj.(*Function)
	if !ok {
		return nil, nil
	}

	fctx := newFunctionContext(ctx, f.Ret)
	for i, p := range f.Params {
		name := ast.Params[i].Name
		if name != nil {
			if !fctx.AddValue(name.Value, p) {
				return nil, utils.Errorf(ast.Name.Position, "duplicate identifier")
			}
		}
	}

	bctx, body, err := analyseBlock(fctx, *ast.Body, false)
	if err != nil {
		return nil, err
	} else if !bctx.IsEnd() {
		if f.Ret.Equal(None) {
			body.Stmts = append(body.Stmts, &Return{})
			bctx.SetEnd()
		} else {
			return nil, utils.Errorf(ast.Name.Position, "function missing return")
		}
	}
	f.Body = body
	return f, nil
}

// 全局变量
func analyseGlobalVariable(ctx *packageContext, astAttrs []parse.Attr, ast parse.GlobalVariable) (*GlobalVariable, utils.Error) {
	if ast.Variable.Type == nil && ast.Variable.Value == nil {
		return nil, utils.Errorf(ast.Variable.Name.Position, "expect a type or a value")
	}

	var typ Type
	var err utils.Error
	if ast.Variable.Type != nil {
		typ, err = analyseType(ctx, ast.Variable.Type)
		if err != nil {
			return nil, err
		}
	}

	var value Expr
	if ast.Variable.Type != nil && ast.Variable.Value != nil {
		value, err = analyseConstantExpr(typ, *ast.Variable.Value)
		if err != nil {
			return nil, err
		}
		value, err = expectExpr(ast.Variable.Value.Position, typ, value)
		if err != nil {
			return nil, err
		}
	} else if ast.Variable.Type == nil && ast.Variable.Value != nil {
		value, err = analyseConstantExpr(nil, *ast.Variable.Value)
		if err != nil {
			return nil, err
		} else if IsNoneType(value.GetType()) {
			return nil, utils.Errorf(ast.Variable.Value.Position, "expect a value")
		}
		typ = value.GetType()
	}

	v := &GlobalVariable{
		Type:  typ,
		Value: value,
	}
	if !ctx.AddValue(ast.Public != nil, ast.Variable.Name.Value, v) {
		return nil, utils.Errorf(ast.Variable.Name.Position, "duplicate identifier")
	}

	// 属性
	var errors []utils.Error
	for _, astAttr := range astAttrs {
		switch {
		case astAttr.Extern != nil:
			v.ExternName = astAttr.Extern.Value
		case astAttr.LinkAsm != nil:
			linkPath := stlos.Path(*astAttr.LinkAsm)
			if !linkPath.IsAbsolute() {
				linkPath = ctx.path.Join(linkPath)
			}
			if !linkPath.IsExist() {
				errors = append(errors, utils.Errorf(astAttr.Position, "can not find path `%s`", linkPath))
			}
			ctx.f.Links[linkPath] = struct{}{}
		case astAttr.LinkLib != nil:
			ctx.f.Libs[string(*astAttr.LinkLib)] = struct{}{}
		case astAttr.NoReturn != nil:
			fallthrough
		case astAttr.Exit != nil:
			errors = append(errors, utils.Errorf(astAttr.Position, "attribute `@noreturn` cannot be used for global variables"))
		default:
			panic("")
		}
	}
	if v.ExternName == "" && ast.Variable.Value == nil {
		errors = append(errors, utils.Errorf(ast.Variable.Name.Position, "missing value"))
	}
	if len(errors) == 0 {
		return v, nil
	} else if len(errors) == 1 {
		return nil, errors[0]
	} else {
		return nil, utils.NewMultiError(errors...)
	}
}
