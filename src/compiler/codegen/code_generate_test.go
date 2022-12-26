package codegen

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"testing"
	"tinygo.org/x/go-llvm"
)

func TestCodegen(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	module := NewCodeGenerator().Codegen(*mean)
	fmt.Println(module)
}

func TestOptimize(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	module := Optimize(NewCodeGenerator().Codegen(*mean), llvm.OptLevelAggressive, llvm.SizeLevelZ)
	fmt.Println(module)
}
