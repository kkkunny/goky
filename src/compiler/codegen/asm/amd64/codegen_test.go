package amd64

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/generate_ssa"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"testing"
)

func TestCodegen(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	ssa := generate_ssa.Optimize(generate_ssa.NewGenerator().Generate(*mean))
	NewCodeGenerator(os.Stdout, *ssa).Codegen()
}
