package codegen

import (
	"tinygo.org/x/go-llvm"
)

// Optimize 优化
func Optimize(module llvm.Module, optLevel llvm.OptLevel, sizeLevel llvm.SizeLevel) llvm.Module {
	pmb := llvm.NewPassManagerBuilder()
	pmb.SetOptLevel(optLevel)
	pmb.SetSizeLevel(sizeLevel)
	pm := llvm.NewPassManager()
	pmb.Populate(pm)

	pm.Run(module)
	return module
}
