package codegen

import (
	"bytes"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/analyse"
	"io"
	"strings"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	writer    io.Writer
	include   map[string]struct{}
	typedef   bytes.Buffer
	vars      map[analyse.Expr]string
	types     map[string]string
	typeCount uint
}

// NewGenerator 新建代码生成器
func NewGenerator(w io.Writer) *CodeGenerator {
	return &CodeGenerator{
		writer: w,
		vars:   make(map[analyse.Expr]string),
		types:  make(map[string]string),
	}
}

func (self *CodeGenerator) writef(f string, a ...any) {
	_, _ = fmt.Fprintf(self.writer, f, a...)
}

// Generate 中间代码生成
func (self *CodeGenerator) Generate(mean analyse.ProgramContext) {
	root := self.writer
	// 声明
	var decl bytes.Buffer
	self.writer = &decl
	var gc uint
	for _, global := range mean.Globals {
		switch g := global.(type) {
		case *analyse.GlobalVariable:
			if g.ExternName == "" {
				g.ExternName = fmt.Sprintf("__g%d", gc)
				gc++
			}
			self.vars[g] = g.ExternName
			self.writef("%s %s = %s;\n", self.generateType(g.GetType()), g.ExternName, self.generateConstantExpr(g.Value))
		case *analyse.Function:
			if g.ExternName == "" {
				g.ExternName = fmt.Sprintf("__g%d", gc)
				gc++
			}
			self.vars[g] = g.ExternName
			paramTypes := make([]string, len(g.Params))
			for i, p := range g.Params {
				paramTypes[i] = self.generateType(p.GetType())
			}
			self.writef("%s %s(%s);\n", self.generateType(g.Ret), g.ExternName, strings.Join(paramTypes, ", "))
		default:
			panic("")
		}
	}
	// 定义
	var def bytes.Buffer
	self.writer = &def
	for _, global := range mean.Globals {
		switch g := global.(type) {
		case *analyse.GlobalVariable:
		case *analyse.Function:
			if g.Body != nil {
				params := make([]string, len(g.Params))
				for i, p := range g.Params {
					name := fmt.Sprintf("__p%d", i)
					params[i] = fmt.Sprintf("%s %s", self.generateType(p.GetType()), name)
					self.vars[p] = name
				}
				self.writef("%s %s(%s)", self.generateType(g.Ret), g.ExternName, strings.Join(params, ", "))
				var count uint
				self.generateBlock(&count, *g.Body)
			}

			// main
			if g.Main {
				self.writef("int main(){\n")
				self.writef("%s();\n", g.ExternName)
				self.writef("return 0;\n")
				self.writef("}\n")
			}
			// init
			// TODO
			// fini
			// TODO
		default:
			panic("")
		}
	}
	// 输出
	self.writer = root
	for s := range self.include {
		self.writef("include <%s>\n", s)
	}
	self.writef("\n")
	_, _ = self.typedef.WriteTo(self.writer)
	self.writef("\n\n")
	_, _ = decl.WriteTo(self.writer)
	self.writef("\n\n")
	_, _ = def.WriteTo(self.writer)
}
