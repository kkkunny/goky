package ir

import (
	"fmt"
	"strings"
)

// Module 模块
type Module struct {
	NamedFunctions   map[string]*Function
	UnNamedFunctions []*Function
	NamedGlobals     map[string]*Global
	UnNamedGlobals   []*Global
}

// NewModule 新建模块
func NewModule() *Module {
	return &Module{
		NamedFunctions: make(map[string]*Function),
		NamedGlobals:   make(map[string]*Global),
	}
}

func (self Module) String() string {
	var buf strings.Builder
	for _, v := range self.NamedGlobals {
		buf.WriteString(v.String())
		buf.WriteByte('\n')
	}
	for _, v := range self.UnNamedGlobals {
		buf.WriteString(v.String())
		buf.WriteByte('\n')
	}
	for _, f := range self.NamedFunctions {
		buf.WriteString(f.String())
		buf.WriteByte('\n')
	}
	for _, f := range self.UnNamedFunctions {
		buf.WriteString(f.String())
		buf.WriteByte('\n')
	}
	return buf.String()
}

// NewFunction 新建函数
func (self *Module) NewFunction(ft *TypeFunc, name string) *Function {
	params := make([]*Param, len(ft.Params))
	for i, p := range ft.Params {
		params[i] = &Param{
			No:   uint(i),
			Type: p,
		}
	}
	f := &Function{
		Belong: self,
		Type:   ft,
		Params: params,
		Name:   name,
	}
	if name != "" {
		if _, ok := self.NamedFunctions[name]; ok {
			panic(fmt.Sprintf("重复的函数：%s", name))
		}
		self.NamedFunctions[name] = f
	} else {
		f.Name = fmt.Sprintf("f%d", len(self.UnNamedFunctions))
		self.UnNamedFunctions = append(self.UnNamedFunctions, f)
	}
	return f
}

// NewGlobal 新建全局变量
func (self *Module) NewGlobal(t Type, name string, v Constant) *Global {
	g := &Global{
		Name:  name,
		Type:  t,
		Value: v,
	}
	if name != "" {
		if _, ok := self.NamedGlobals[name]; ok {
			panic(fmt.Sprintf("重复的全局变量：%s", name))
		}
		self.NamedGlobals[name] = g
	} else {
		g.Name = fmt.Sprintf("g%d", len(self.UnNamedGlobals))
		self.UnNamedGlobals = append(self.UnNamedGlobals, g)
	}
	return g
}

// Function 函数
type Function struct {
	Belong   *Module
	Type     *TypeFunc
	Name     string
	Params   []*Param
	Blocks   []*Block
	varCount uint

	// 属性
	NoReturn bool
}

func (self Function) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s %s", self.Type.Ret, self.GetName()))
	buf.WriteByte('(')
	for i, p := range self.Params {
		buf.WriteString(p.Type.String())
		buf.WriteByte(' ')
		buf.WriteString(p.GetName())
		if i < len(self.Params)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteByte(')')
	if self.NoReturn {
		buf.WriteString(" #noreturn")
	}
	if len(self.Blocks) > 0 {
		buf.WriteString(":\n")
		for _, b := range self.Blocks {
			buf.WriteString(b.String())
		}
	} else {
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (self *Function) NewBlock() *Block {
	block := &Block{
		Depends: make(map[*Block]struct{}),
		Belong:  self,
		No:      uint(len(self.Blocks)),
	}
	self.Blocks = append(self.Blocks, block)
	return block
}

func (self Function) GetName() string {
	return self.Name
}

func (self Function) GetType() Type {
	return NewPtrType(self.Type)
}

func (self Function) stmt() {}

func (self Function) GetParam(i uint) *Param {
	return self.Params[i]
}

func (self Function) GetBelong() *Block {
	return nil
}

// Param 参数
type Param struct {
	No   uint
	Type Type
}

func (self Param) stmt() {}

func (self Param) GetName() string {
	return fmt.Sprintf("p%d", self.No)
}

func (self Param) GetType() Type {
	return NewPtrType(self.Type)
}

func (self Param) GetBelong() *Block {
	return nil
}

// Global 全局变量
type Global struct {
	Name  string
	Type  Type
	Value Constant
}

func (self Global) String() string {
	if self.Value != nil {
		return fmt.Sprintf("%s %s = %s", self.Type, self.Name, self.Value.GetName())
	} else {
		return fmt.Sprintf("%s %s", self.Type, self.Name)
	}
}

func (self Global) GetName() string {
	return self.Name
}

func (self Global) GetType() Type {
	return NewPtrType(self.Type)
}

func (self Global) stmt() {}

func (self Global) GetBelong() *Block {
	return nil
}
