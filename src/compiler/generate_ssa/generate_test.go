package generate_ssa

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"os"
	"syscall"
	"testing"
)

func TestGenerate(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	module := NewGenerator().Generate(*mean)
	syscall.Chroot()
	fmt.Println(module)
}

func TestOptimize(t *testing.T) {
	util.Must(os.Setenv("KROOT", stlos.Path(util.MustValue(os.Getwd())).GetParent().GetParent().GetParent().String()))
	ast := util.MustValue(parse.ParseFile("../../../main.k"))
	mean := util.MustValue(analyse.AnalyseMain(*ast))
	module := NewGenerator().Generate(*mean)
	module = Optimize(module)
	fmt.Println(module)
}
