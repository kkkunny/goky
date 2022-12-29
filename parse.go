//go:build test && parse

package main

import (
	"encoding/json"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/parse"
	"github.com/kkkunny/stl/util"
)

func main() {
	ast := util.MustValue(parse.ParseFile("main.k"))
	out := string(util.MustValue(json.MarshalIndent(ast, "", "  ")))
	fmt.Println(out)
}
