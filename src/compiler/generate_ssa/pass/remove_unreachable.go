package pass

import "github.com/kkkunny/klang/src/compiler/ir_ssa"

type RemoveUnreachable struct{}

func (self RemoveUnreachable) Run(module *ir.Module) {
	for _, f := range module.NamedFunctions {
		self.runFunction(f)
	}
	for _, f := range module.UnNamedFunctions {
		self.runFunction(f)
	}
}

func (self RemoveUnreachable) runFunction(function *ir.Function) {
	var blocks []*ir.Block
	for _, b := range function.Blocks {
		self.runBlock(b)
		if len(b.Stmts) != 0 {
			blocks = append(blocks, b)
		}
	}
	function.Blocks = blocks
}

func (self RemoveUnreachable) runBlock(block *ir.Block) {
	if len(block.Stmts) == 0 {
		return
	}
	var end int
loop:
	for i, s := range block.Stmts {
		switch s.(type) {
		case *ir.Return, *ir.CondGoto, *ir.Goto, *ir.Unreachable:
			end = i
			break loop
		}
	}
	block.Stmts = block.Stmts[:end+1]
}
