//go:build test && analyse

package main

import (
	"encoding/json"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/parse"
	"github.com/kkkunny/stl/util"
)

func main() {
	ast := util.MustValue(parse.ParseFile("main.k"))
	mean := util.MustValue(analyse.AnalyseMain(ast))
	out := string(util.MustValue(json.MarshalIndent(mean, "", "  ")))
	fmt.Println(out)
}
