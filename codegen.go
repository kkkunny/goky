//go:build test && codegen

package main

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/codegen"
	"github.com/kkkunny/klang/src/compiler/parse"
	"github.com/kkkunny/stl/util"
)

func main() {
	ast := util.MustValue(parse.ParseFile("main.k"))
	mean := util.MustValue(analyse.AnalyseMain(ast))
	module := codegen.NewCodeGenerator().Codegen(*mean)
	fmt.Println(module)
}
