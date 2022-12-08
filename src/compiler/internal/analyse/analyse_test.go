package analyse

import (
	"encoding/json"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"testing"
)

func TestAnalyse(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../../main.k"))
	mean := util.MustValue(AnalyseMain(*ast))
	out := string(util.MustValue(json.MarshalIndent(mean, "", "  ")))
	fmt.Println(out)
}
