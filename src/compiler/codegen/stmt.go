package codegen

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
)

// 代码块
func (self *CodeGenerator) generateBlock(mean analyse.Block) {
	self.writef("{\n")
	for _, stmt := range mean.Stmts {
		self.generateStmt(stmt)
	}
	self.writef("}\n")
}

// 语句
func (self *CodeGenerator) generateStmt(mean analyse.Stmt) {
	switch stmt := mean.(type) {
	case *analyse.Return:
		if stmt.Value == nil {
			self.doBeforeFuncEnd()
			self.writef("return;\n")
		} else {
			ret := self.generateExpr(stmt.Value)
			self.doBeforeFuncEnd()
			self.writef("return %s;\n", ret)
		}
	case *analyse.Variable:
		name := fmt.Sprintf("__v%d", self.varCount)
		self.varCount++
		self.writef("%s %s = %s;\n", self.generateType(stmt.GetType()), name, self.generateExpr(stmt.Value))
		self.vars[stmt] = name
	case analyse.Expr:
		self.writef("%s;\n", self.generateExpr(stmt))
	case *analyse.Block:
		self.generateBlock(*stmt)
	case *analyse.IfElse:
		self.writef("if(%s)", self.generateExpr(stmt.Cond))
		self.generateBlock(*stmt.True)
		self.writef("else")
		self.generateBlock(*stmt.False)
	case *analyse.Loop:
		self.writef("while(%s)", self.generateExpr(stmt.Cond))
		self.generateBlock(*stmt.Body)
	case *analyse.LoopControl:
		self.writef("%s;\n", stmt.Type)
	case *analyse.Defer:
		self.defers.Add(stmt.Call)
	default:
		panic("")
	}
}

func (self *CodeGenerator) doBeforeFuncEnd() {
	for iter := self.defers.Iterator(); iter.HasValue(); iter.Next() {
		self.writef(fmt.Sprintf("%s;\n", self.generateExpr(iter.Value())))
	}
}
