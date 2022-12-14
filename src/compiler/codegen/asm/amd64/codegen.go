package amd64

import (
	"fmt"
	ir "github.com/kkkunny/klang/src/compiler/ir_ssa"
	"github.com/kkkunny/klang/src/compiler/utils"
	"io"
)

var argReg8 = [...]string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argReg16 = [...]string{"di", "si", "dx", "cx", "r8w", "r9w"}
var argReg32 = [...]string{"edi", "esi", "edx", "ecx", "r8d", "r9d"}
var argReg64 = [...]string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var argReg128 = [...]string{"xmm0", "xmm1", "xmm2", "xmm3", "xmm4", "xmm5", "xmm6", "xmm7"}

// CodeGenerator 代码生成器
type CodeGenerator struct {
	writer io.Writer
	ssa    ir.Module

	stackDepth int
	values     map[ir.Value]int
	zeros      map[uint]struct{}
}

// NewCodeGenerator 新建代码生成器
func NewCodeGenerator(w io.Writer, ssa ir.Module) *CodeGenerator {
	return &CodeGenerator{
		writer: w,
		ssa:    ssa,

		values: make(map[ir.Value]int),
	}
}

// 写入
func (self *CodeGenerator) write(f string, a ...any) {
	_, _ = fmt.Fprintf(self.writer, f, a...)
}

// 写入
func (self *CodeGenerator) writeln(f string, a ...any) {
	self.write(f, a...)
	self.write("\n")
}

// 入栈
func (self *CodeGenerator) push(r string) {
	self.writeln("  push %s", r)
	self.stackDepth++
}

// 出栈
func (self *CodeGenerator) pop(r string) {
	self.writeln("  pop %s", r)
	self.stackDepth--
}

// 入栈
func (self *CodeGenerator) push128(r string) {
	self.writeln("  sub rsp, 8")
	self.writeln("  movsd qword ptr [rsp], %s", r)
	self.stackDepth++
}

// 出栈
func (self *CodeGenerator) pop128(r string) {
	self.writeln("  movsd %s, qword ptr [rsp]", r)
	self.writeln("  add rsp, 8")
	self.stackDepth--
}

// 入栈
func (self *CodeGenerator) pushBigSize(size uint, r string) {
	realSize := utils.AlignTo(size, 8)
	self.writeln("  sub rsp, %d", realSize)
	for i := uint(0); i < size; i++ {
		self.writeln("  mov sil, byte ptr [%s + %d]", r, i)
		self.writeln("  mov byte ptr [rsp + %d], sil", i)
	}
	self.stackDepth += int(realSize) / 8
}

// load
func (self *CodeGenerator) load(t ir.Type) {
	switch typ := t.(type) {
	case ir.IntType, *ir.TypePtr:
		switch {
		case typ.GetByte() == 1:
			self.writeln("  mov al, byte ptr [rax]")
		case typ.GetByte() == 2:
			self.writeln("  mov ax, word ptr [rax]")
		case typ.GetByte() <= 4:
			self.writeln("  mov eax, dword ptr [rax]")
		case typ.GetByte() <= 8:
			self.writeln("  mov rax, qword ptr [rax]")
		default:
			panic("")
		}
	case *ir.TypeFloat:
		switch typ.GetByte() {
		case 4:
			self.writeln("  movss xmm0, dword ptr [rax]")
		case 8:
			self.writeln("  movsd xmm0, qword ptr [rax]")
		default:
			panic("")
		}
	case *ir.TypeArray:
	default:
		panic("")
	}
}

// store
func (self *CodeGenerator) store(t ir.Type) {
	switch typ := t.(type) {
	case ir.IntType, *ir.TypePtr:
		switch {
		case t.GetByte() == 1:
			self.writeln("  mov byte ptr [rdi], al")
		case t.GetByte() == 2:
			self.writeln("  mov word ptr [rdi], ax")
		case t.GetByte() <= 4:
			self.writeln("  mov dword ptr [rdi], eax")
		case t.GetByte() <= 8:
			self.writeln("  mov qword ptr [rdi], rax")
		default:
			panic("")
		}
	case *ir.TypeFloat:
		switch typ.GetByte() {
		case 4:
			self.writeln("  movss dword ptr [rdi], xmm0")
		case 8:
			self.writeln("  movsd qword ptr [rdi], xmm0")
		default:
			panic("")
		}
	case *ir.TypeArray:
		for i := uint(0); i < typ.GetByte(); i++ {
			self.writeln("  mov sil, byte ptr [rax + %d]", i)
			self.writeln("  mov byte ptr [rdi + %d], sil", i)
		}
	default:
		panic("")
	}
}

// Codegen 代码生成
func (self *CodeGenerator) Codegen() {
	for i, ssa := range self.ssa.UnNamedGlobals {
		ssa.Name = fmt.Sprintf("__unnamed_%d", i+1)
	}
	for i, ssa := range self.ssa.UnNamedFunctions {
		ssa.Name = fmt.Sprintf("__unnamed_%d", len(self.ssa.UnNamedGlobals)+i+1)
	}

	self.writeln("  .intel_syntax noprefix")
	for _, ssa := range self.ssa.NamedFunctions {
		self.codegenFunction(true, *ssa)
	}
	for _, ssa := range self.ssa.UnNamedFunctions {
		self.codegenFunction(false, *ssa)
	}
	for _, ssa := range self.ssa.NamedGlobals {
		self.codegenGlobal(true, *ssa)
	}
	for _, ssa := range self.ssa.UnNamedGlobals {
		self.codegenGlobal(false, *ssa)
	}

	for size := range self.zeros {
		self.codegenGlobalZero(size)
	}
}
