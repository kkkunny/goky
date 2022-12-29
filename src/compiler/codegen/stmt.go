package codegen

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/klang/src/compiler/analyse"
)

// 代码块
func (self *CodeGenerator) codegenBlock(mean analyse.Block) {
	for _, stmt := range mean.Stmts {
		self.codegenStmt(stmt)
	}
}

// 语句
func (self *CodeGenerator) codegenStmt(mean analyse.Stmt) {
	switch meanStmt := mean.(type) {
	case *analyse.Return:
		self.codegenReturn(*meanStmt)
	case *analyse.Variable:
		self.codegenVariable(meanStmt)
	case analyse.Expr:
		self.codegenExpr(meanStmt, true)
	case *analyse.Block:
		self.codegenBlock(*meanStmt)
	case *analyse.IfElse:
		self.codegenIfElse(*meanStmt)
	case *analyse.Loop:
		self.codegenLoop(*meanStmt)
	case *analyse.LoopControl:
		self.codegenLoopControl(*meanStmt)
	case *analyse.Defer:
		self.codegenDefer(*meanStmt)
	default:
		panic("")
	}
}

// 函数返回
func (self *CodeGenerator) codegenReturn(mean analyse.Return) {
	if mean.Value == nil {
		self.doneBeforeFuncEnd()
		self.builder.CreateRetVoid()
	} else {
		value := self.codegenExpr(mean.Value, true)
		self.doneBeforeFuncEnd()
		self.builder.CreateRet(value)
	}
}

// 变量
func (self *CodeGenerator) codegenVariable(mean *analyse.Variable) {
	alloca := self.builder.CreateAlloca(self.codegenType(mean.Type), "")
	value := self.codegenExpr(mean.Value, true)
	self.vars[mean] = alloca
	self.builder.CreateStore(value, alloca)
}

// 条件分支
func (self *CodeGenerator) codegenIfElse(mean analyse.IfElse) {
	cond := self.codegenExpr(mean.Cond, true)
	tb := llvm.AddBasicBlock(self.function, "")
	if mean.False == nil {
		eb := llvm.AddBasicBlock(self.function, "")
		self.builder.CreateCondBr(cond, tb, eb)

		self.builder.SetInsertPointAtEnd(tb)
		self.codegenBlock(*mean.True)
		self.builder.CreateBr(eb)

		self.builder.SetInsertPointAtEnd(eb)
	} else {
		fb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
		self.builder.CreateCondBr(cond, tb, fb)

		self.builder.SetInsertPointAtEnd(tb)
		self.codegenBlock(*mean.True)
		self.builder.CreateBr(eb)

		self.builder.SetInsertPointAtEnd(fb)
		self.codegenBlock(*mean.False)
		self.builder.CreateBr(eb)

		self.builder.SetInsertPointAtEnd(eb)
	}
}

// 循环
func (self *CodeGenerator) codegenLoop(mean analyse.Loop) {
	cb := llvm.AddBasicBlock(self.function, "")
	self.builder.CreateBr(cb)

	self.builder.SetInsertPointAtEnd(cb)
	lb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
	self.builder.CreateCondBr(self.codegenExpr(mean.Cond, true), lb, eb)

	cbBk, ebBk := self.cb, self.eb
	self.cb, self.eb = cb, eb
	self.builder.SetInsertPointAtEnd(lb)
	self.codegenBlock(*mean.Body)
	self.builder.CreateBr(cb)
	self.cb, self.eb = cbBk, ebBk

	self.builder.SetInsertPointAtEnd(eb)
}

// 循环控制
func (self *CodeGenerator) codegenLoopControl(mean analyse.LoopControl) {
	if mean.Type == "break" {
		self.builder.CreateBr(self.eb)
	} else {
		self.builder.CreateBr(self.cb)
	}
}

type deferInfo struct {
	Func llvm.Value
	Args []llvm.Value
}

// 延迟调用
func (self *CodeGenerator) codegenDefer(mean analyse.Defer) {
	f := self.codegenExpr(mean.Call.Func, true)
	args := make([]llvm.Value, len(mean.Call.Args))
	for i, a := range mean.Call.Args {
		args[i] = self.codegenExpr(a, true)
	}
	self.defers = append(self.defers, deferInfo{
		Func: f,
		Args: args,
	})
}

// 函数结束前需做的事情
func (self *CodeGenerator) doneBeforeFuncEnd() {
	for _, d := range self.defers {
		self.builder.CreateCall(d.Func.Type().ReturnType(), d.Func, d.Args, "")
	}
}
