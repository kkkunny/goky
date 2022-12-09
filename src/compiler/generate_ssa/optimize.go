package generate_ssa

import (
	"github.com/kkkunny/klang/src/compiler/generate_ssa/pass"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
)

var passes = [...]pass.Pass{
	new(pass.RemoveUnreachable),
}

// Optimize 优化
func Optimize(module *ir.Module) *ir.Module {
	for _, p := range passes {
		p.Run(module)
	}
	return module
}
