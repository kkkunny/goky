package codegen

import (
	"github.com/kkkunny/klang/src/compiler/internal/analyse"
	"github.com/kkkunny/klang/src/compiler/internal/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"testing"
)

func TestCodegen(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	NewGenerator(os.Stdout).Generate(*mean)
}
