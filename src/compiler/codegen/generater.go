package codegen

import (
	"bytes"
	"fmt"
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/stl/list"
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
	varCount  uint
	defers    *list.SingleLinkedList[*analyse.Call]
}

// NewGenerator 新建代码生成器
func NewGenerator(w io.Writer) *CodeGenerator {
	return &CodeGenerator{
		writer: w,
		vars:   make(map[analyse.Expr]string),
		types:  make(map[string]string),
		defers: list.NewSingleLinkedList[*analyse.Call](),
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
				self.varCount = 0
				self.defers.Clear()
				self.generateBlock(*g.Body)
			}
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
