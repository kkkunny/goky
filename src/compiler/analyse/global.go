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

// 函数模板
type functionTemplate struct {
	ast      *parse.FunctionHead
	astAttrs []parse.Attr
	impls    map[string]*Function
}

func (self functionTemplate) global() {}

func (self functionTemplate) stmt() {}

func (self functionTemplate) ident() {}

func (self functionTemplate) GetType() Type {
	return None
}

func (self functionTemplate) GetMut() bool {
	return false
}

func (self functionTemplate) IsTemporary() bool {
	return true
}

// *********************************************************************************************************************

// 函数声明
func analyseFunctionDecl(ctx *packageContext, astAttrs []parse.Attr, ast parse.FunctionHead) (*Function, utils.Error) {
	retType, err := analyseType(ctx, ast.Tail.Function.Ret)
	if err != nil {
		return nil, err
	}

	params := make([]*Param, len(ast.Tail.Function.Params))
	var errors []utils.Error
	for i, p := range ast.Tail.Function.Params {
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
			errors = append(errors, utils.Errorf(astAttr.Position, "attribute cannot be used for global variables"))
		}
	}
	if f.ExternName == "" && ast.Tail.Function.Body == nil {
		errors = append(errors, utils.Errorf(ast.Tail.Function.Name.Position, "missing function body"))
	}
	if len(errors) == 1 {
		return nil, errors[0]
	} else if len(errors) > 1 {
		return nil, utils.NewMultiError(errors...)
	}

	if !ctx.AddValue(ast.Public != nil, ast.Tail.Function.Name.Value, f) {
		return nil, utils.Errorf(ast.Tail.Function.Name.Position, "duplicate identifier")
	}
	return f, nil
}

// 函数定义
func analyseFunctionDef(ctx *packageContext, f *Function, ast parse.Function) utils.Error {
	fctx := newFunctionContext(ctx, f.Ret)
	for i, p := range f.Params {
		name := ast.Params[i].Name
		if name != nil {
			if !fctx.AddValue(name.Value, p) {
				return utils.Errorf(name.Position, "duplicate identifier")
			}
		}
	}

	bctx, body, err := analyseBlock(fctx, *ast.Body, false)
	if err != nil {
		return err
	} else if !bctx.IsEnd() {
		if f.Ret.Equal(None) {
			body.Stmts = append(body.Stmts, &Return{})
			bctx.SetEnd()
		} else {
			return utils.Errorf(ast.Name.Position, "function missing return")
		}
	}
	f.Body = body
	return nil
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
		default:
			errors = append(errors, utils.Errorf(astAttr.Position, "attribute cannot be used for global variables"))
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

// 类型定义
func analyseTypedef(ctx *packageContext, ast parse.Typedef) utils.Error {
	dst, err := analyseType(ctx, &ast.Dst)
	if err != nil {
		return err
	}
	td := ctx.typedefs[ast.Name.Value].Second
	td.Dst = dst
	return nil
}

// 方法声明
func analyseMethodDecl(ctx *packageContext, astAttrs []parse.Attr, ast parse.FunctionHead) (*Function, utils.Error) {
	_selfType, err := analyseType(ctx, &parse.Type{Ident: &parse.TypeIdent{
		Position: ast.Tail.Method.Self.Position,
		Name:     ast.Tail.Method.Self,
	}})
	if err != nil {
		return nil, err
	}
	selfType := NewPtrType(_selfType)

	retType, err := analyseType(ctx, ast.Tail.Method.Ret)
	if err != nil {
		return nil, err
	}

	params := make([]*Param, len(ast.Tail.Method.Params)+1)
	params[0] = &Param{Type: selfType}
	var errors []utils.Error
	for i, p := range ast.Tail.Method.Params {
		pt, err := analyseType(ctx, &p.Type)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		params[i+1] = &Param{Type: pt}
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
		case astAttr.NoReturn != nil:
			f.NoReturn = true
		case astAttr.Exit != nil:
			f.Exit = true
			f.NoReturn = true
		default:
			errors = append(errors, utils.Errorf(astAttr.Position, "attribute cannot be used for global variables"))
		}
	}
	if len(errors) == 1 {
		return nil, errors[0]
	} else if len(errors) > 1 {
		return nil, utils.NewMultiError(errors...)
	}

	name := _selfType.String() + "." + ast.Tail.Method.Name.Value
	if !ctx.AddValue(ast.Public != nil, name, f) {
		return nil, utils.Errorf(ast.Tail.Method.Name.Position, "duplicate identifier")
	}
	return f, nil
}

// 方法定义
func analyseMethodDef(ctx *packageContext, ast parse.Method) utils.Error {
	_selfType, err := analyseType(ctx, &parse.Type{Ident: &parse.TypeIdent{
		Position: ast.Self.Position,
		Name:     ast.Self,
	}})
	if err != nil {
		return err
	}

	name := _selfType.String() + "." + ast.Name.Value
	f := ctx.GetValue(name).Second.(*Function)
	fctx := newFunctionContext(ctx, f.Ret)
	for i, p := range f.Params {
		if i == 0 {
			fctx.AddValue("self", p)
		} else {
			pn := ast.Params[i].Name
			if pn != nil {
				if !fctx.AddValue(pn.Value, p) {
					return utils.Errorf(pn.Position, "duplicate identifier")
				}
			}
		}
	}

	bctx, body, err := analyseBlock(fctx, ast.Body, false)
	if err != nil {
		return err
	} else if !bctx.IsEnd() {
		if f.Ret.Equal(None) {
			body.Stmts = append(body.Stmts, &Return{})
			bctx.SetEnd()
		} else {
			return utils.Errorf(ast.Name.Position, "function missing return")
		}
	}
	f.Body = body
	return nil
}

// 函数模板声明
func analyseFunctionTemplateDecl(ctx *packageContext, astAttrs []parse.Attr, ast *parse.FunctionHead) utils.Error {
	// 属性
	var errors []utils.Error
	for _, astAttr := range astAttrs {
		switch {
		case astAttr.NoReturn != nil:
		case astAttr.Exit != nil:
		default:
			errors = append(errors, utils.Errorf(astAttr.Position, "attribute cannot be used for global variables"))
		}
	}
	if ast.Tail.Function.Body == nil {
		errors = append(errors, utils.Errorf(ast.Tail.Function.Name.Position, "missing function body"))
	}
	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 1 {
		return utils.NewMultiError(errors...)
	}

	ft := &functionTemplate{
		ast:      ast,
		astAttrs: astAttrs,
		impls:    make(map[string]*Function),
	}
	if !ctx.AddValue(ast.Public != nil, ast.Tail.Function.Name.Value, ft) {
		return utils.Errorf(ast.Tail.Function.Name.Position, "duplicate identifier")
	}

	return nil
}
