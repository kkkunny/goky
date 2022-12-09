package parse

import (
	"encoding/json"
	"fmt"
	"github.com/kkkunny/stl/util"
	"testing"
)

func TestParse(t *testing.T) {
	ast := util.MustValue(ParseFile("../../../main.k"))
	out := string(util.MustValue(json.MarshalIndent(ast, "", "  ")))
	fmt.Println(out)
}
