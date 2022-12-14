package amd64

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
	"github.com/kkkunny/klang/src/compiler/utils"
	stlutil "github.com/kkkunny/stl/util"
	"math"
)

// 函数
func (self *CodeGenerator) codegenFunction(pub bool, ssa ir.Function) {
	if len(ssa.Blocks) == 0 {
		return
	}

	// 函数头
	self.writeln("  .text")
	if pub {
		self.writeln("  .globl %s", ssa.Name)
	} else {
		self.writeln("  .local %s", ssa.Name)
	}
	self.writeln("  .p2align 4, 0x90")
	self.writeln("  .type %s,@function", ssa.Name)
	self.writeln("%s:", ssa.Name)
	stackSize, regParams := self.getFuncStackSize(ssa)
	self.writeln("  enter %d, 0", stackSize)
	if ssa.Name == "main" {
		fmt.Println("")
	}
	// 参数
	self.pushParams(ssa.Type.Ret, regParams)
	// 函数体
	for _, b := range ssa.Blocks {
		self.codegenBlock(b)
	}
	if self.stackDepth != 0 {
		panic("")
	}
	// 函数返回
	self.writeln(".L.%s.return:", ssa.Name)
	self.writeln("  leave")
	self.writeln("  ret")
}

// 获取函数栈大小和寄存器传递参数并设置值位置
func (self *CodeGenerator) getFuncStackSize(f ir.Function) (res uint, params []*ir.Param) {
	// 返回值
	var retStack bool
	if !ir.IsBasicType(f.Type.Ret) {
		if f.Type.Ret.GetByte() > 16 {
			retStack = true
			res += 8
		}
	}
	// 区分寄存器和栈传递参数并计算参数栈大小
	gp, fp := stlutil.Ternary(retStack, 1, 0), 0
	var stack uint = 16
	regParams, stackParams := make([]*ir.Param, 0, len(f.Params)), make([]*ir.Param, 0, len(f.Params))
	for _, p := range f.Params {
		gpbk, fpbk := gp, fp
		if ir.IsFloatType(p.Type) && fp < len(argReg128) {
			fp++
		} else if !ir.IsBasicType(p.Type) {
			if p.Type.GetByte() <= 16 {
				f, s := self.analyseTypeRegister(0, p.Type)
				if p.Type.GetByte() <= 8 {
					if f && gp < len(argReg64) {
						gp++
					} else if !f && fp < len(argReg128) {
						fp++
					}
				} else {
					if f && s && gp < len(argReg64)-1 {
						gp += 2
					} else if f && !s && gp < len(argReg64) && fp < len(argReg128) {
						gp++
						fp++
					} else if s && !f && gp < len(argReg64) && fp < len(argReg128) {
						gp++
						fp++
					} else if !f && !s && fp < len(argReg128)-1 {
						fp += 2
					}
				}
			}
		} else if gp < len(argReg64) {
			gp++
		}
		if gpbk != gp || fpbk != fp {
			regParams = append(regParams, p)
			res += uint((gp - gpbk) * 8)
			res += uint((fp - fpbk) * 8)
			self.values[p] = -int(res)
		} else {
			stackParams = append(stackParams, p)
			self.values[p] = int(stack)
			if !ir.IsBasicType(p.Type) {
				stack += utils.AlignTo(p.Type.GetByte(), 8)
			} else {
				stack += 8
			}
		}
	}
	// 本地变量
	for _, b := range f.Blocks {
		for _, s := range b.Stmts {
			switch stmt := s.(type) {
			case *ir.Alloc:
				res += stmt.Type.GetByte()
				self.values[stmt] = -int(res)
			case *ir.Call:
				if stmt.No != 0 {
					res += 8
					self.values[stmt] = -int(res)
				}
			case ir.Value:
				res += 8
				self.values[stmt] = -int(res)
			}
		}
	}
	// 对齐16字节
	res = utils.AlignTo(res, 16)
	return res, regParams
}

// 入栈寄存器传递参数
func (self *CodeGenerator) pushParams(ret ir.Type, params []*ir.Param) {
	// 返回值
	var retStack bool
	if !ir.IsBasicType(ret) {
		if ret.GetByte() > 16 {
			self.writeln("  mov qword ptr [rbp - 8], rdi")
			retStack = true
		}
	}
	// 参数
	gp, fp := stlutil.Ternary(retStack, 1, 0), 0
	for _, p := range params {
		if ir.IsFloatType(p.Type) && fp < len(argReg128) {
			switch {
			case p.Type.GetByte() <= 4:
				self.writeln("  movss dword ptr [rbp - %d], %s", -self.values[p], argReg128[fp])
			case p.Type.GetByte() <= 8:
				self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p], argReg128[fp])
			default:
				panic("")
			}
			fp++
		} else if !ir.IsBasicType(p.Type) {
			f, s := self.analyseTypeRegister(0, p.Type)
			if p.Type.GetByte() <= 8 {
				if f && gp < len(argReg64) {
					self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p], argReg64[gp])
					gp++
				} else if !f && fp < len(argReg128) {
					self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p], argReg128[fp])
					fp++
				}
			} else {
				if f && s && gp < len(argReg64)-1 {
					self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p], argReg64[gp])
					self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p]-8, argReg64[gp+1])
					gp += 2
				} else if f && !s && gp < len(argReg64) && fp < len(argReg128) {
					self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p], argReg64[gp])
					gp++
					self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p]-8, argReg128[fp])
					fp++
				} else if s && !f && gp < len(argReg64) && fp < len(argReg128) {
					self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p], argReg128[fp])
					fp++
					self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p]-8, argReg64[gp])
					gp++
				} else if !f && !s && fp < len(argReg128)-1 {
					self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p], argReg128[fp])
					self.writeln("  movsd qword ptr [rbp - %d], %s", -self.values[p]-8, argReg128[fp+1])
					fp += 2
				}
			}
		} else if gp < len(argReg64) {
			switch {
			case p.Type.GetByte() == 1:
				self.writeln("  mov byte ptr [rbp - %d], %s", -self.values[p], argReg8[gp])
			case p.Type.GetByte() == 2:
				self.writeln("  mov word ptr [rbp - %d], %s", -self.values[p], argReg16[gp])
			case p.Type.GetByte() <= 4:
				self.writeln("  mov dword ptr [rbp - %d], %s", -self.values[p], argReg32[gp])
			case p.Type.GetByte() <= 8:
				self.writeln("  mov qword ptr [rbp - %d], %s", -self.values[p], argReg64[gp])
			default:
				panic("")
			}
			gp++
		}
	}
}

// 全局变量
func (self *CodeGenerator) codegenGlobal(pub bool, ssa ir.Global) {
	if ssa.Value == nil {
		return
	}

	self.writeln("  .type %s,@object", ssa.Name)
	if !ssa.Value.IsZero() {
		self.writeln("  .data")
	} else {
		self.writeln("  .bss")
	}
	if pub {
		self.writeln("  .globl %s", ssa.Name)
	} else {
		self.writeln("  .local %s", ssa.Name)
	}
	self.writeln("  .p2align %d", int(math.Log2(float64(ssa.Type.GetByte()))))
	self.writeln("%s:", ssa.Name)
	self.codegenConstant(ssa.Value)
	self.writeln("  .size %s, %d", ssa.Name, ssa.Type.GetByte())
}

// 空变量
func (self *CodeGenerator) codegenGlobalZero(size uint) {
	name := fmt.Sprintf("__zero_%d", size)
	self.writeln("  .type %s,@object", name)
	self.writeln("  .bss")
	self.writeln("  .local %s", name)
	self.writeln("%s:", name)
	self.writeln("  .zero %d", size)
	self.writeln("  .size %s, %d", name, size)
}
