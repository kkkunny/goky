//go:build test && parse

package main

import (
	"encoding/json"
	"fmt"
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
	out := string(util.MustValue(json.MarshalIndent(ast, "", "  ")))
	fmt.Println(out)
}
