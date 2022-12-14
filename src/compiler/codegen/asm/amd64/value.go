package amd64

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
	"math"
)

// 值
func (self *CodeGenerator) codegenValue(ssa ir.Value) {
	switch value := ssa.(type) {
	case *ir.Int:
		self.writeln("  mov rax, %d", value.Value)
	case *ir.Float:
		switch value.GetType().GetByte() {
		case 4:
			self.writeln("  mov eax, %d", math.Float32bits(float32(value.Value)))
		case 8:
			self.writeln("  movabs rax, %d", math.Float64bits(value.Value))
		default:
			panic("")
		}
		self.writeln("  movq xmm0, rax")
	case *ir.Empty:
		switch t := value.GetType().(type) {
		case *ir.TypePtr:
			self.writeln("  mov rax, 0")
		case *ir.TypeArray, *ir.TypeStruct:
			size := t.GetByte()
			self.zeros[size] = struct{}{}
			name := fmt.Sprintf("__zero_%d", size)
			self.writeln("  lea rax, [rip + %s]", name)
		default:
			panic("")
		}
	case *ir.Function:
		if len(value.Blocks) != 0 {
			self.writeln("  lea rax, [rip + %s]", value.Name)
		} else {
			self.writeln("  mov rax, [rip + %s@GOTPCREL]", value.Name)
		}
	case *ir.Param:
		self.writeln("  lea rax, [rbp - %d]", -self.values[ssa])
	case *ir.Global:
		self.writeln("  lea rax, [rip + %s]", value.Name)
	case *ir.Alloc:
		self.writeln("  lea rax, [rbp - %d]", -self.values[ssa])
	case ir.Stmt:
		if ir.IsFloatType(ssa.GetType()) {
			switch ssa.GetType().GetByte() {
			case 4:
				self.writeln("  movss xmm0, qword ptr [rbp - %d]", -self.values[ssa])
			case 8:
				self.writeln("  movsd xmm0, qword ptr [rbp - %d]", -self.values[ssa])
			default:
				panic("")
			}
		} else {
			self.writeln("  mov rax, qword ptr [rbp - %d]", -self.values[ssa])
		}
	default:
		panic("")
	}
}

// 常量
func (self *CodeGenerator) codegenConstant(ssa ir.Constant) {
	if ssa.IsZero() {
		self.writeln("  .zero %d", ssa.GetType().GetByte())
	}
	switch value := ssa.(type) {
	case *ir.Int:
		switch ssa.GetType().GetByte() {
		case 1:
			self.writeln("  .byte %d", value.Value)
		case 2:
			self.writeln("  .short %d", value.Value)
		case 4:
			self.writeln("  .long %d", value.Value)
		case 8:
			self.writeln("  .quad %d", value.Value)
		default:
			panic("")
		}
	case *ir.Float:
		switch ssa.GetType().GetByte() {
		case 4:
			self.writeln("  .long %d", math.Float32bits(float32(value.Value)))
		case 8:
			self.writeln("  .quad %d", math.Float64bits(value.Value))
		default:
			panic("")
		}
	case *ir.Empty:
	case *ir.Array:
		for _, e := range value.Elems {
			self.codegenConstant(e)
		}
	case *ir.Struct:
		t := value.GetType().(*ir.TypeStruct)
		for i, e := range value.Elems {
			offset := t.GetOffset(uint(i))
			self.codegenConstant(e)
			var next uint
			if i < len(value.Elems)-1 {
				next = t.GetOffset(uint(i) + 1)
			} else {
				next = t.GetByte()
			}
			diff := next - offset - e.GetType().GetByte()
			if diff > 0 {
				self.writeln("  .zero %d", diff)
			}
		}
	default:
		panic("")
	}
}
