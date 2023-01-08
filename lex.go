//go:build test && lex

package main

import (
	"fmt"
	"github.com/kkkunny/Sim/src/compiler/lex"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"strings"
)

func main() {
	k := stlos.Path(os.Args[1])
	lexer := lex.NewLexer(k, strings.NewReader(string(util.MustValue(os.ReadFile(k.String())))))
	for tok := lexer.Scan(); tok.Kind != lex.EOF; tok = lexer.Scan() {
		fmt.Println(tok)
	}
}
