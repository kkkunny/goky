//go:build test && codegen

package main

import (
	"fmt"
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/codegen"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
)

func main() {
	path := stlos.Path(os.Args[1])
	var ast *parse.Package
	if path.IsDir() {
		ast = util.MustValue(parse.ParsePackage(stlos.Path(os.Args[1])))
	} else {
		ast = util.MustValue(parse.ParseFile(stlos.Path(os.Args[1])))
	}
	mean := util.MustValue(analyse.AnalyseMain(ast))
	module := codegen.NewCodeGenerator().Codegen(*mean)
	util.Must(llvm.VerifyModule(module, llvm.ReturnStatusAction))
	fmt.Println(module)
}
