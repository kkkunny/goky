package analyse

import (
	"github.com/kkkunny/klang/src/compiler/internal/parse"
	"github.com/kkkunny/klang/src/compiler/internal/utils"
)

// Block 代码块
type Block struct {
	Stmts []Stmt
}

func (self Block) stmt() {}

// Stmt 语句
type Stmt interface {
	stmt()
}

// Return 函数返回
type Return struct {
	Value Expr
}

func (self Return) stmt() {}

// Variable 变量
type Variable struct {
	Type  Type
	Value Expr
}

func (self Variable) stmt() {}

func (self Variable) ident() {}

func (self Variable) GetType() Type {
	return self.Type
}

func (self Variable) GetMut() bool {
	return true
}

func (self Variable) IsTemporary() bool {
	return false
}

// IfElse 条件分支
type IfElse struct {
	Cond        Expr
	True, False *Block
}

func (self IfElse) stmt() {}

// Loop 循环
type Loop struct {
	Cond Expr
	Body *Block
}

func (self Loop) stmt() {}

// LoopControl 循环
type LoopControl struct {
	Type string
}

func (self LoopControl) stmt() {}

// Defer 延迟调用
type Defer struct {
	Call *Call
}

func (self Defer) stmt() {}

// *********************************************************************************************************************

// 代码块
func analyseBlock(ctx localContext, ast parse.Block, inLoop bool) (*blockContext, *Block, utils.Error) {
	bctx := newBlockContext(ctx, inLoop)

	var stmts []Stmt

	var errors []utils.Error
	for _, s := range ast.Stmts {
		if bctx.IsEnd() {
			break
		}
		stmt, err := analyseStmt(bctx, s)
		if err != nil {
			errors = append(errors, err)
		} else {
			stmts = append(stmts, stmt)
		}
	}

	block := &Block{Stmts: stmts}
	if len(errors) == 0 {
		return bctx, block, nil
	} else if len(errors) == 1 {
		return nil, nil, errors[0]
	} else {
		return nil, nil, utils.NewMultiError(errors...)
	}
}

// 语句
func analyseStmt(ctx *blockContext, ast parse.Stmt) (Stmt, utils.Error) {
	switch {
	case ast.Return != nil:
		stmt, err := analyseReturn(ctx, *ast.Return)
		ctx.SetEnd()
		return stmt, err
	case ast.Variable != nil:
		return analyseVariable(ctx, *ast.Variable)
	case ast.Expr != nil:
		return analyseExpr(ctx, nil, *ast.Expr)
	case ast.Block != nil:
		bctx, stmt, err := analyseBlock(ctx, *ast.Block, false)
		if err != nil {
			return nil, err
		}
		if bctx.IsEnd() {
			ctx.SetEnd()
		}
		return stmt, nil
	case ast.IfElse != nil:
		stmt, err, end := analyseIfElse(ctx, *ast.IfElse)
		if err != nil {
			return nil, err
		}
		if end {
			ctx.SetEnd()
		}
		return stmt, nil
	case ast.For != nil:
		return analyseFor(ctx, *ast.For)
	case ast.LoopControl != nil:
		if !ctx.IsInLoop() {
			return nil, utils.Errorf(ast.Position, "must in a loop")
		}
		ctx.SetEnd()
		return &LoopControl{Type: *ast.LoopControl}, nil
	case ast.Defer != nil:
		return analyseDefer(ctx, ast)
	default:
		panic("")
	}
}

// 函数返回
func analyseReturn(ctx *blockContext, ast parse.Return) (*Return, utils.Error) {
	ret := ctx.GetRetType()
	if ast.Value == nil {
		if ret.Equal(None) {
			return &Return{}, nil
		} else {
			return nil, utils.Errorf(ast.Position, "expect a return value")
		}
	} else {
		if ret.Equal(None) {
			return nil, utils.Errorf(ast.Position, "not expect a return value")
		} else {
			value, err := analyseExpr(ctx, ret, *ast.Value)
			if err != nil {
				return nil, err
			}
			value, err = expectExpr(ast.Value.Position, ret, value)
			return &Return{Value: value}, err
		}
	}
}

// 变量
func analyseVariable(ctx *blockContext, ast parse.Variable) (*Variable, utils.Error) {
	if ast.Type == nil && ast.Value == nil {
		return nil, utils.Errorf(ast.Name.Position, "expect a type or a value")
	}

	var typ Type
	var err utils.Error
	if ast.Type != nil {
		typ, err = analyseType(ctx.GetPackageContext(), ast.Type)
		if err != nil {
			return nil, err
		}
	}

	var value Expr
	if ast.Type != nil && ast.Value != nil {
		value, err = analyseExpr(ctx, typ, *ast.Value)
		if err != nil {
			return nil, err
		}
		value, err = expectExpr(ast.Value.Position, typ, value)
		if err != nil {
			return nil, err
		}
	} else if ast.Type == nil && ast.Value != nil {
		value, err = analyseExpr(ctx, nil, *ast.Value)
		if err != nil {
			return nil, err
		} else if IsNoneType(value.GetType()) {
			return nil, utils.Errorf(ast.Value.Position, "expect a value")
		}
		typ = value.GetType()
	} else {
		value = getDefaultExprByType(typ)
	}

	v := &Variable{
		Type:  typ,
		Value: value,
	}
	if !ctx.AddValue(ast.Name.Value, v) {
		return nil, utils.Errorf(ast.Name.Position, "duplicate identifier")
	}
	return v, nil
}

// 条件分支
func analyseIfElse(ctx *blockContext, ast parse.IfElse) (*IfElse, utils.Error, bool) {
	cond, err := analyseExpr(ctx, Bool, ast.Cond)
	if err != nil {
		return nil, err, false
	}
	cond, err = expectExpr(ast.Cond.Position, Bool, cond)
	if err != nil {
		return nil, err, false
	}

	tctx, tb, te := analyseBlock(ctx, ast.True, false)

	if ast.Suffix == nil {
		if te != nil {
			return nil, te, false
		}
		return &IfElse{
			Cond: cond,
			True: tb,
		}, nil, false
	} else if ast.Suffix.False != nil {
		fctx, fb, fe := analyseBlock(ctx, *ast.Suffix.False, false)
		if te != nil && fe != nil {
			return nil, utils.NewMultiError(te, fe), false
		} else if te != nil {
			return nil, te, false
		} else if fe != nil {
			return nil, fe, false
		}
		return &IfElse{
			Cond:  cond,
			True:  tb,
			False: fb,
		}, nil, tctx.IsEnd() && fctx.IsEnd()
	} else {
		nb, ne, nret := analyseIfElse(ctx, *ast.Suffix.Else)
		if te != nil && ne != nil {
			return nil, utils.NewMultiError(te, ne), false
		} else if te != nil {
			return nil, te, false
		} else if ne != nil {
			return nil, ne, false
		}
		return &IfElse{
			Cond:  cond,
			True:  tb,
			False: &Block{Stmts: []Stmt{nb}},
		}, nil, tctx.IsEnd() && nret
	}
}

// 循环
func analyseFor(ctx *blockContext, ast parse.For) (*Loop, utils.Error) {
	cond, err := analyseExpr(ctx, Bool, ast.Cond)
	if err != nil {
		return nil, err
	}
	cond, err = expectExpr(ast.Cond.Position, Bool, cond)
	if err != nil {
		return nil, err
	}

	_, body, err := analyseBlock(ctx, ast.Body, true)
	if err != nil {
		return nil, err
	}

	return &Loop{
		Cond: cond,
		Body: body,
	}, nil
}

// 延迟调用
func analyseDefer(ctx *blockContext, ast parse.Stmt) (*Defer, utils.Error) {
	obj, err := analysePrimaryPostfix(ctx, nil, *ast.Defer)
	if err != nil {
		return nil, err
	}
	call, ok := obj.(*Call)
	if !ok {
		return nil, utils.Errorf(ast.Position, "expect a function call")
	}
	return &Defer{Call: call}, nil
}
