package generate_ssa

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
)

// 代码块
func (self *Generator) generateBlock(mean analyse.Block) {
	for _, stmt := range mean.Stmts {
		self.generateStmt(stmt)
	}
}

// 语句
func (self *Generator) generateStmt(mean analyse.Stmt) {
	switch meanStmt := mean.(type) {
	case *analyse.Return:
		self.generateReturn(*meanStmt)
	case *analyse.Variable:
		self.generateVariable(meanStmt)
	case analyse.Expr:
		self.generateExpr(meanStmt, true)
	case *analyse.Block:
		self.generateBlock(*meanStmt)
	case *analyse.IfElse:
		self.generateIfElse(*meanStmt)
	case *analyse.Loop:
		self.generateLoop(*meanStmt)
	case *analyse.LoopControl:
		self.generateLoopControl(*meanStmt)
	case *analyse.Defer:
		self.generateDefer(*meanStmt)
	default:
		panic("")
	}
}

// 函数返回
func (self *Generator) generateReturn(mean analyse.Return) {
	if mean.Value == nil {
		self.doneBeforeFuncEnd()
		self.block.NewReturn(nil)
	} else {
		value := self.generateExpr(mean.Value, true)
		self.doneBeforeFuncEnd()
		self.block.NewReturn(value)
	}
}

// 变量
func (self *Generator) generateVariable(mean *analyse.Variable) {
	typ := generateType(mean.Type)
	alloca := self.block.NewAlloc(typ)
	value := self.generateExpr(mean.Value, true)
	self.vars[mean] = alloca
	self.block.NewStore(value, alloca)
}

// 条件分支
func (self *Generator) generateIfElse(mean analyse.IfElse) {
	cond := self.generateExpr(mean.Cond, true)
	tb := self.block.Belong.NewBlock()
	if mean.False == nil {
		pb := self.block
		self.block = tb
		self.generateBlock(*mean.True)
		tb = self.block

		eb := self.block.Belong.NewBlock()
		pb.NewCondGoto(cond, tb, eb)
		tb.NewGoto(eb)
		self.block = eb
	} else {
		fb := self.block.Belong.NewBlock()
		self.block.NewCondGoto(cond, tb, fb)

		self.block = tb
		self.generateBlock(*mean.True)
		tb = self.block

		self.block = fb
		self.generateBlock(*mean.False)
		fb = self.block

		eb := self.block.Belong.NewBlock()
		tb.NewGoto(eb)
		fb.NewGoto(eb)

		self.block = eb
	}
}

// 循环
func (self *Generator) generateLoop(mean analyse.Loop) {
	cb := self.block.Belong.NewBlock()
	self.block.NewGoto(cb)

	self.block = cb
	cond := self.generateExpr(mean.Cond, true)
	lb, eb := self.block.Belong.NewBlock(), self.block.Belong.NewBlock()
	self.block.NewCondGoto(cond, lb, eb)

	cbBk, ebBk := self.cb, self.eb
	self.cb, self.eb = cb, eb
	self.block = lb
	self.generateBlock(*mean.Body)
	self.block.NewGoto(cb)
	self.cb, self.eb = cbBk, ebBk

	self.block = eb
}

// 循环控制
func (self *Generator) generateLoopControl(mean analyse.LoopControl) {
	if mean.Type == "break" {
		self.block.NewGoto(self.eb)
	} else {
		self.block.NewGoto(self.cb)
	}
}

type deferInfo struct {
	Func ir.Value
	Args []ir.Value
}

// 延迟调用
func (self *Generator) generateDefer(mean analyse.Defer) {
	f := self.generateExpr(mean.Call.Func, true)
	args := make([]ir.Value, len(mean.Call.Args))
	for i, a := range mean.Call.Args {
		args[i] = self.generateExpr(a, true)
	}
	self.defers = append(self.defers, deferInfo{
		Func: f,
		Args: args,
	})
}

// 函数结束前需做的事情
func (self *Generator) doneBeforeFuncEnd() {
	// defer
	for _, d := range self.defers {
		self.block.NewCall(d.Func, d.Args...)
	}
}
