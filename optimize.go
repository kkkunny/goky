//go:build test && optimize

package main

import (
	"fmt"
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/codegen"
	"github.com/kkkunny/stl/util"
)

func main() {
	ast := util.MustValue(parse.ParseFile("main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	module := codegen.Optimize(codegen.NewCodeGenerator().Codegen(*mean), llvm.OptLevelAggressive, llvm.SizeLevelZ)
	fmt.Println(module)
}
