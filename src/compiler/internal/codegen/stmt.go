package codegen

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/analyse"
)

// 代码块
func (self *CodeGenerator) generateBlock(varCount *uint, mean analyse.Block) {
	self.writef("{\n")
	for _, stmt := range mean.Stmts {
		self.generateStmt(varCount, stmt)
	}
	self.writef("}\n")
}

// 语句
func (self *CodeGenerator) generateStmt(varCount *uint, mean analyse.Stmt) {
	switch stmt := mean.(type) {
	case *analyse.Return:
		if stmt.Value == nil {
			self.writef("return;\n")
		} else {
			self.writef("return %s;\n", self.generateExpr(stmt.Value))
		}
	case *analyse.Variable:
		name := fmt.Sprintf("__v%d", *varCount)
		*varCount += 1
		self.writef("%s %s = %s;\n", self.generateType(stmt.GetType()), name, self.generateExpr(stmt.Value))
		self.vars[stmt] = name
	case analyse.Expr:
		self.writef("%s;\n", self.generateExpr(stmt))
	case *analyse.Block:
		self.generateBlock(varCount, *stmt)
	case *analyse.IfElse:
		// TODO
		panic("")
	case *analyse.Loop:
		// TODO
		panic("")
	case *analyse.LoopControl:
		self.writef("%s;\n", stmt.Type)
	case *analyse.Defer:
		// TODO
		panic("")
	default:
		panic("")
	}
}
