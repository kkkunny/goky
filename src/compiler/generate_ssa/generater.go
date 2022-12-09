package generate_ssa

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
)

// Generator 中间代码生成器
type Generator struct {
	module *ir.Module
	block  *ir.Block

	vars map[analyse.Expr]ir.Value

	// loop
	cb, eb *ir.Block
	// defer
	defers []deferInfo
}

// NewGenerator 新建中间代码生成器
func NewGenerator() *Generator {
	return &Generator{
		module: ir.NewModule(),
		vars:   make(map[analyse.Expr]ir.Value),
	}
}

// Generate 中间代码生成
func (self *Generator) Generate(mean analyse.ProgramContext) *ir.Module {
	// 声明
	for _, g := range mean.Globals {
		switch global := g.(type) {
		case *analyse.Function:
			ft := generateType(global.GetType()).(*ir.TypePtr).Elem.(*ir.TypeFunc)
			f := self.module.NewFunction(ft, global.ExternName)
			if global.NoReturn || global.Exit {
				f.NoReturn = true
			}
			self.vars[global] = f
		case *analyse.GlobalVariable:
			vt := generateType(global.GetType())
			self.vars[global] = self.module.NewGlobal(vt, global.ExternName, nil)
		default:
			panic("")
		}
	}
	// 定义
	for _, g := range mean.Globals {
		switch global := g.(type) {
		case *analyse.Function:
			if global.Body != nil {
				f := self.vars[global].(*ir.Function)
				for i, p := range global.Params {
					self.vars[p] = f.GetParam(uint(i))
				}
				self.block = f.NewBlock()
				self.generateBlock(*global.Body)
				self.defers = nil
			}
		case *analyse.GlobalVariable:
			if global.Value != nil {
				self.vars[global].(*ir.Global).Value = self.generateConstantExpr(global.Value)
			}
		default:
			panic("")
		}
	}
	return self.module
}
