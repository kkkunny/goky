package pass

import "github.com/kkkunny/klang/src/compiler/ir_ssa"

type Pass interface {
	Run(module *ir.Module)
}
