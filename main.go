package main

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/codegen"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
)

func main() {
	from := stlos.Path(os.Args[1])
	from = util.MustValue(from.GetAbsolute())

	var ast *parse.Package
	if from.IsDir() {
		ast = util.MustValue(parse.ParsePackage(from))
	} else {
		ast = util.MustValue(parse.ParseFile(from))
	}

	mean := util.MustValue(analyse.AnalyseMain(*ast))

	cfilepath := from.WithExtension("c")
	cfile := util.MustValue(os.OpenFile(cfilepath.String(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755))
	codegen.NewGenerator(cfile).Generate(*mean)
	util.Must(cfile.Close())
	defer os.Remove(cfilepath.String())

	outfilepath := from.WithExtension("out")
	util.Must(stlos.Exec("gcc", "-fPIC", "-O2", "-o", outfilepath.String(), cfilepath.String()))
}
