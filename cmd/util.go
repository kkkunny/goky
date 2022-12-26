package cmd

import (
	"errors"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/codegen"
	"github.com/kkkunny/klang/src/compiler/parse"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"tinygo.org/x/go-llvm"
)

// LookupCmd 查找命令
func LookupCmd(cmd ...string) (string, *exec.Cmd) {
	for _, c := range cmd {
		p, err := exec.LookPath(c)
		if err != nil {
			continue
		}
		e := exec.Command(p)
		e.Stdin = os.Stdin
		e.Stdout = os.Stdout
		e.Stderr = os.Stderr
		return c, e
	}
	return "", nil
}

// RandomString 随机字符串
func RandomString(n uint8) string {
	rand.Seed(time.Now().Unix())
	var buf strings.Builder
	for i := uint8(0); i < n; i++ {
		n := rand.Intn(62)
		if n < 26 {
			buf.WriteByte('a' + byte(n))
		} else if n < 52 {
			buf.WriteByte('A' + byte(n) - 26)
		} else {
			buf.WriteByte('0' + byte(n) - 52)
		}
	}
	return buf.String()
}

// 输出汇编
func outputAsm(config *buildConfig, from, to stlos.Path) (stlos.Path, error) {
	if to == "" {
		for {
			to = stlos.Path(os.TempDir()).Join(stlos.Path(RandomString(6) + ".s"))
			if !to.IsExist() {
				break
			}
		}
	}

	var ast *parse.Package
	var err error
	if from.IsDir() {
		ast, err = parse.ParsePackage(from)
	} else {
		ast, err = parse.ParseFile(from)
	}
	if err != nil {
		return "", err
	}
	mean, err := analyse.AnalyseMain(*ast)
	if err != nil {
		return "", err
	}
	module := codegen.Optimize(
		codegen.NewCodeGenerator().Codegen(*mean),
		util.Ternary(config.Release, llvm.OptLevelAggressive, llvm.OptLevelDefault),
		util.Ternary(config.Release, llvm.SizeLevelZ, llvm.SizeLevelNone),
	)

	if err = llvm.InitializeNativeTarget(); err != nil {
		return "", err
	}
	if err = llvm.InitializeNativeAsmPrinter(); err != nil {
		return "", err
	}
	module.SetTarget(llvm.DefaultTargetTriple())
	target, err := llvm.GetTargetFromTriple(module.Target())
	if err != nil {
		return "", err
	}
	tm := target.CreateTargetMachine(module.Target(), "generic", "", llvm.CodeGenLevelNone, llvm.RelocPIC, llvm.CodeModelDefault)
	if err = tm.EmitToFile(module, to.String(), llvm.AssemblyFile); err != nil {
		return "", err
	}

	for l := range mean.Links {
		config.Linkages = append(config.Linkages, l)
	}
	for l := range mean.Libs {
		config.Libraries = append(config.Libraries, l)
	}
	return to, nil
}

// 输出目标文件
func outputObject(from, to stlos.Path, links []stlos.Path) (stlos.Path, error) {
	if to == "" {
		for {
			to = stlos.Path(os.TempDir()).Join(stlos.Path(RandomString(6) + ".o"))
			if !to.IsExist() {
				break
			}
		}
	}

	_, assembler := LookupCmd("as")
	if assembler == nil {
		return "", errors.New("can not found a assembler")
	}

	assembler.Args = append(assembler.Args, "-o", to.String(), from.String())
	for _, link := range links {
		assembler.Args = append(assembler.Args, link.String())
	}
	return to, assembler.Run()
}

// 输出动态库文件
func outputSharedFile(from, to stlos.Path, libraries, libraryPaths []string) (stlos.Path, error) {
	if to == "" {
		for {
			to = stlos.Path(os.TempDir()).Join("lib" + stlos.Path(RandomString(6)) + ".so")
			if !to.IsExist() {
				break
			}
		}
	}

	_, linker := LookupCmd("clang", "gcc")
	if linker == nil {
		return "", errors.New("can not found a linker")
	}
	linker.Args = append(linker.Args, "-shared", "-fPIC", "-o", to.String(), from.String())
	for _, l := range libraries {
		linker.Args = append(linker.Args, fmt.Sprintf("-l%s", l))
	}
	for _, L := range libraryPaths {
		linker.Args = append(linker.Args, fmt.Sprintf("-L%s", L))
	}
	return to, linker.Run()
}

// 输出可执行文件
func outputExecutableFile(from, to stlos.Path, libraries, libraryPaths []string) (stlos.Path, error) {
	if to == "" {
		for {
			to = stlos.Path(os.TempDir()).Join(stlos.Path(RandomString(6)))
			if !to.IsExist() {
				break
			}
		}
	}

	_, linker := LookupCmd("clang", "gcc")
	if linker == nil {
		return "", errors.New("can not found a linker")
	}
	linker.Args = append(linker.Args, "-static", "-fPIC", "-o", to.String(), from.String())
	for _, l := range libraries {
		linker.Args = append(linker.Args, fmt.Sprintf("-l%s", l))
	}
	for _, L := range libraryPaths {
		linker.Args = append(linker.Args, fmt.Sprintf("-L%s", L))
	}
	return to, linker.Run()
}
