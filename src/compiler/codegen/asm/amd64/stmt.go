package amd64

import (
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
	"github.com/kkkunny/klang/src/compiler/utils"
	stlutil "github.com/kkkunny/stl/util"
)

// 代码块
func (self *CodeGenerator) codegenBlock(ssa *ir.Block) {
	self.writeln(".L.%s.block.%d:", ssa.Belong.Name, ssa.No)
	for _, s := range ssa.Stmts {
		self.codegenStmt(s)
	}
}

// 语句
func (self *CodeGenerator) codegenStmt(ssa ir.Stmt) {
	switch stmt := ssa.(type) {
	case *ir.Return:
		if stmt.Value != nil {
			self.codegenValue(stmt.Value)
			self.pushRet(stmt.Value.GetType())
		}
		self.writeln("  jmp .L.%s.return", stmt.Belong.Belong.Name)
	case *ir.Add:
		if ir.IsIntType(stmt.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  add rax, rdi")
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  addss xmm0, xmm1")
			case 8:
				self.writeln("  addsd xmm0, xmm1")
			default:
				panic("")
			}
		}
	case *ir.Sub:
		if ir.IsIntType(stmt.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  sub rax, rdi")
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  subss xmm0, xmm1")
			case 8:
				self.writeln("  subsd xmm0, xmm1")
			default:
				panic("")
			}
		}
	case *ir.Mul:
		if ir.IsIntType(stmt.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  imul rax, rdi")
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  mulss xmm0, xmm1")
			case 8:
				self.writeln("  mulsd xmm0, xmm1")
			default:
				panic("")
			}
		}
	case *ir.Div:
		if ir.IsIntType(stmt.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  %s", stlutil.Ternary(stmt.GetType().GetByte() == 8, "cqo", "cdq"))
				self.writeln("  idiv rdi")
			} else {
				self.writeln("  xor ecx, ecx")
				self.writeln("  mov edx, ecx")
				self.writeln("  div rdi")
			}
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  divss xmm0, xmm1")
			case 8:
				self.writeln("  divsd xmm0, xmm1")
			default:
				panic("")
			}
		}
	case *ir.Mod:
		if ir.IsIntType(stmt.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  %s", stlutil.Ternary(stmt.GetType().GetByte() == 8, "cqo", "cdq"))
				self.writeln("  idiv rdi")
			} else {
				self.writeln("  xor ecx, ecx")
				self.writeln("  mov edx, ecx")
				self.writeln("  div rdi")
			}
			self.writeln("  mov rax, rdx")
		} else {
			// TODO
			panic("")
		}
	case *ir.And:
		self.codegenValue(stmt.Right)
		self.push("rax")
		self.codegenValue(stmt.Left)
		self.pop("rdi")
		self.writeln("  and rax, rdi")
	case *ir.Or:
		self.codegenValue(stmt.Right)
		self.push("rax")
		self.codegenValue(stmt.Left)
		self.pop("rdi")
		self.writeln("  or rax, rdi")
	case *ir.Xor:
		self.codegenValue(stmt.Right)
		self.push("rax")
		self.codegenValue(stmt.Left)
		self.pop("rdi")
		self.writeln("  xor rax, rdi")
	case *ir.Shl:
		self.codegenValue(stmt.Right)
		self.push("rax")
		self.codegenValue(stmt.Left)
		self.pop("rcx")
		self.writeln("  shl rax, cl")
	case *ir.Shr:
		self.codegenValue(stmt.Right)
		self.push("rax")
		self.codegenValue(stmt.Left)
		self.pop("rcx")
		if ir.IsSintType(stmt.GetType()) {
			self.writeln("  sar rax, cl")
		} else {
			self.writeln("  shr rax, cl")
		}
	case *ir.Alloc:
		return
	case *ir.Store:
		self.codegenValue(stmt.To)
		self.push("rax")
		self.codegenValue(stmt.From)
		self.pop("rdi")
		self.store(stmt.From.GetType())
	case *ir.Load:
		self.codegenValue(stmt.Ptr)
		self.load(stmt.GetType())
	case *ir.Call:
		// 入参
		stack := self.pushArgs(stmt)
		// 调用
		if f, ok := stmt.Func.(*ir.Function); ok {
			self.writeln("  call %s%s", f.Name, stlutil.Ternary(len(f.Blocks) != 0, "", "@PLT"))
		} else {
			self.codegenValue(stmt.Func)
			self.writeln("  call rax")
		}
		self.writeln("  add rsp, %d", stack)
		if stmt.No == 0 {
			return
		}
		// 返回值
		self.popRet(stmt.GetType())
	case *ir.ArrayIndex:
		self.codegenValue(stmt.Index)
		self.push("rax")
		self.codegenValue(stmt.From)
		self.pop("rdi")
		self.writeln("  imul rdi, %d", stmt.GetElemType().GetByte())
		self.writeln("  add rax, rdi")
	case *ir.PtrIndex:
		self.codegenValue(stmt.Index)
		self.push("rax")
		self.codegenValue(stmt.From)
		self.pop("rdi")
		self.writeln("  imul rdi, %d", stmt.GetElemType().GetByte())
		self.writeln("  add rax, rdi")
	case *ir.StructIndex:
		self.codegenValue(stmt.From)
		self.writeln("  add rax, %d", stmt.GetFromType().GetOffset(stmt.Index))
	case *ir.Eq:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			self.writeln("  sete al")
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  sete al")
			self.writeln("  setnp cl")
			self.writeln("  and al, cl")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Ne:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			self.writeln("  setne al")
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  setne al")
			self.writeln("  setp cl")
			self.writeln("  or al, cl")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Lt:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  setl al")
			} else {
				self.writeln("  setb al")
			}
		} else {
			self.codegenValue(stmt.Left)
			self.push128("xmm0")
			self.codegenValue(stmt.Right)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  seta al")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Le:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  setle al")
			} else {
				self.writeln("  setbe al")
			}
		} else {
			self.codegenValue(stmt.Left)
			self.push128("xmm0")
			self.codegenValue(stmt.Right)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  setae al")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Gt:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  setg al")
			} else {
				self.writeln("  seta al")
			}
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  seta al")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Ge:
		if ir.IsIntType(stmt.Left.GetType()) {
			self.codegenValue(stmt.Right)
			self.push("rax")
			self.codegenValue(stmt.Left)
			self.pop("rdi")
			self.writeln("  cmp rax, rdi")
			if ir.IsSintType(stmt.GetType()) {
				self.writeln("  setge al")
			} else {
				self.writeln("  setae al")
			}
		} else {
			self.codegenValue(stmt.Right)
			self.push128("xmm0")
			self.codegenValue(stmt.Left)
			self.pop128("xmm1")
			switch stmt.GetType().GetByte() {
			case 4:
				self.writeln("  ucomiss xmm0, xmm1")
			case 8:
				self.writeln("  ucomisd xmm0, xmm1")
			default:
				panic("")
			}
			self.writeln("  seta al")
		}
		self.writeln("  and al, 1")
		self.writeln("  movzx eax, al")
	case *ir.Goto:
		self.writeln("  mov r12, %d", stmt.Belong.No)
		self.writeln("  jmp .L.%s.block.%d", stmt.Target.Belong.Name, stmt.Target.No)
	case *ir.CondGoto:
		self.codegenValue(stmt.Cond)
		self.writeln("  mov r12, %d", stmt.Belong.No)
		self.writeln("  cmp rax, 0")
		self.writeln("  je .L.%s.block.%d", stmt.False.Belong.Name, stmt.False.No)
		self.writeln("  jmp .L.%s.block.%d", stmt.True.Belong.Name, stmt.True.No)
	case *ir.Phi:
		for _, b := range stmt.Froms {
			self.writeln("  cmp r12, %d", b.No)
			self.writeln("  je .L.%s.block.%d.phi.from.%d", stmt.Belong.Belong.Name, stmt.Belong.No, b.No)
		}
		for i, b := range stmt.Froms {
			self.writeln(".L.%s.block.%d.phi.from.%d:", stmt.Belong.Belong.Name, stmt.Belong.No, b.No)
			self.codegenValue(stmt.Values[i])
			self.writeln("  jmp .L.%s.block.%d.phi.end", stmt.Belong.Belong.Name, stmt.Belong.No)
		}
		self.writeln(".L.%s.block.%d.phi.end:", stmt.Belong.Belong.Name, stmt.Belong.No)
	case *ir.Itoi:
		self.codegenValue(stmt.From)
		ft := stmt.From.GetType()
		if ft.GetByte() < stmt.To.GetByte() {
			diff := (stmt.To.GetByte() - ft.GetByte()) * 8
			self.writeln("  shl rax, %d", diff)
			self.writeln("  %s rax, %d", stlutil.Ternary(ir.IsSintType(ft), "sar", "shr"), diff)
		}
	case *ir.Ftof:
		self.codegenValue(stmt.From)
		if stmt.From.GetType().GetByte() != stmt.To.GetByte() {
			self.writeln("  %s xmm0, xmm0", stlutil.Ternary(stmt.To.GetByte() == 4, "cvtss2sd", "cvtsd2ss"))
		}
	case *ir.Itof:
		self.codegenValue(stmt.From)
		ft := stmt.From.GetType()
		if ft.GetByte() < stmt.To.GetByte() {
			diff := (stmt.To.GetByte() - ft.GetByte()) * 8
			self.writeln("  shl rax, %d", diff)
			self.writeln("  %s rax, %d", stlutil.Ternary(ir.IsSintType(ft), "sar", "shr"), diff)
		}
		self.writeln("  %s xmm0, eax", stlutil.Ternary(stmt.To.GetByte() == 4, "cvtsi2ss", "cvtsi2sd"))
	case *ir.Ftoi:
		self.codegenValue(stmt.From)
		ft := stmt.From.GetType()
		self.writeln("  %s eax, xmm0", stlutil.Ternary(ft.GetByte() == 4, "cvttss2si", "cvttsd2si"))
		if stmt.To.GetByte() > ft.GetByte() {
			diff := (stmt.To.GetByte() - ft.GetByte()) * 8
			self.writeln("  shl rax, %d", diff)
			self.writeln("  %s rax, %d", stlutil.Ternary(ir.IsSintType(stmt.To), "sar", "shr"), diff)
		}
	case *ir.Ptop:
		self.codegenValue(stmt.From)
	case *ir.Itop:
		self.codegenValue(stmt.From)
		ft := stmt.From.GetType()
		if ft.GetByte() < stmt.To.GetByte() {
			diff := (stmt.To.GetByte() - ft.GetByte()) * 8
			self.writeln("  shl rax, %d", diff)
			self.writeln("  shr rax, %d", diff)
		}
	case *ir.Ptoi:
		self.codegenValue(stmt.From)
		ft := stmt.From.GetType()
		if ft.GetByte() < stmt.To.GetByte() {
			diff := (stmt.To.GetByte() - ft.GetByte()) * 8
			if ir.IsSintType(stmt.To) {
				self.writeln("  shl rax, %d", diff)
				self.writeln("  sar rax, %d", diff)
			} else {
				self.writeln("  shl rax, %d", diff)
				self.writeln("  shr rax, %d", diff)
			}
		}
	case *ir.Unreachable:
		return
	default:
		panic("")
	}

	if value, ok := ssa.(ir.Value); ok {
		if ir.IsFloatType(value.GetType()) {
			switch value.GetType().GetByte() {
			case 4:
				self.writeln("  movss qword ptr [rbp - %d], xmm0", -self.values[value])
			case 8:
				self.writeln("  movsd qword ptr [rbp - %d], xmm0", -self.values[value])
			default:
				panic("")
			}
		} else {
			self.writeln("  mov qword ptr [rbp - %d], rax", -self.values[value])
		}
	}
}

// 入参并返回调用后需要出栈的字节数
func (self *CodeGenerator) pushArgs(call *ir.Call) (res uint) {
	// 返回值
	ret := call.GetType()
	var retStack bool
	if !ir.IsBasicType(ret) {
		if ret.GetByte() > 16 {
			retStack = true
			self.writeln("  sub rsp, %d", utils.AlignTo(ret.GetByte(), 16))
			self.writeln("  lea rdi, [rsp]")
		}
	}
	// 参数入栈
	gp, fp := stlutil.Ternary(retStack, 1, 0), 0
	regArgs, stackArgs := make([]ir.Value, 0, len(call.Args)), make([]ir.Value, 0, len(call.Args))
	for _, a := range call.Args {
		at := a.GetType()
		gpbk, fpbk := gp, fp
		if ir.IsFloatType(at) && fp < len(argReg128) {
			fp++
		} else if !ir.IsBasicType(at) {
			if at.GetByte() <= 16 {
				f, s := self.analyseTypeRegister(0, at)
				if at.GetByte() <= 8 {
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
			regArgs = append(regArgs, a)
		} else {
			stackArgs = append(stackArgs, a)
			res += utils.AlignTo(at.GetByte(), 8)
		}
	}
	// 对齐16字节
	allStack := uint(self.stackDepth*8) + res
	nextStack := utils.AlignTo(allStack, 16)
	diffStack := nextStack - allStack
	if diffStack > 0 {
		self.writeln("  sub rsp, %d", diffStack)
		res += diffStack
	}
	// 先入栈栈传递参数
	for i := len(stackArgs) - 1; i >= 0; i-- {
		arg := stackArgs[i]
		self.codegenValue(arg)
		at := arg.GetType()
		if ir.IsFloatType(at) {
			self.push128("xmm0")
		} else if !ir.IsBasicType(at) {
			self.pushBigSize(at.GetByte(), "rax")
		} else {
			self.push("rax")
		}
		self.stackDepth -= int(utils.AlignTo(at.GetByte(), 8)) / 8
	}
	// 再入栈寄存器传递参数
	for i := len(regArgs) - 1; i >= 0; i-- {
		arg := regArgs[i]
		self.codegenValue(arg)
		at := arg.GetType()
		if ir.IsFloatType(at) {
			self.push128("xmm0")
		} else if !ir.IsBasicType(at) {
			self.pushBigSize(at.GetByte(), "rax")
		} else {
			self.push("rax")
		}
	}
	// 移到寄存器
	gp, fp = stlutil.Ternary(retStack, 1, 0), 0
	for _, a := range regArgs {
		at := a.GetType()
		if ir.IsFloatType(at) && fp < len(argReg128) {
			self.pop128(argReg128[fp])
			fp++
		} else if !ir.IsBasicType(at) {
			f, s := self.analyseTypeRegister(0, at)
			if at.GetByte() <= 8 {
				if f && gp < len(argReg64) {
					self.pop(argReg64[gp])
					gp++
				} else if !f && fp < len(argReg128) {
					self.pop128(argReg128[fp])
					fp++
				}
			} else {
				if f && s && gp < len(argReg64)-1 {
					self.pop(argReg64[gp])
					self.pop(argReg64[gp+1])
					gp += 2
				} else if f && !s && gp < len(argReg64) && fp < len(argReg128) {
					self.pop(argReg64[gp])
					gp++
					self.pop128(argReg128[fp])
					fp++
				} else if s && !f && gp < len(argReg64) && fp < len(argReg128) {
					self.pop128(argReg128[fp])
					fp++
					self.pop(argReg64[gp])
					gp++
				} else if !f && !s && fp < len(argReg128)-1 {
					self.pop128(argReg128[fp])
					self.pop128(argReg128[fp+1])
					fp += 2
				}
			}
		} else if gp < len(argReg64) {
			self.pop(argReg64[gp])
			gp++
		}
	}
	return res
}

// 函数返回值入栈
func (self *CodeGenerator) pushRet(t ir.Type) {
	if !ir.IsBasicType(t) {
		if t.GetByte() <= 16 {
			f, s := self.analyseTypeRegister(0, t)
			if t.GetByte() <= 4 {
				if f {
					self.writeln("  mov eax, dword ptr [rax]")
				} else {
					self.writeln("  movss xmm0, dword ptr [rax]")
				}
			} else if t.GetByte() <= 8 {
				if f {
					self.writeln("  mov rax, qword ptr [rax]")
				} else {
					self.writeln("  movsd xmm0, qword ptr [rax]")
				}
			} else {
				if f && s {
					self.writeln("  mov rdx, rax")
					self.writeln("  mov rax, qword ptr [rdx]")
					self.writeln("  mov rdx, qword ptr [rdx + 8]")
				} else if f {
					self.writeln("  movsd xmm0, qword ptr [rax + 8]")
					self.writeln("  mov rax, qword ptr [rax]")
				} else if s {
					self.writeln("  movsd xmm0, qword ptr [rax]")
					self.writeln("  mov rax, qword ptr [rax + 8]")
				} else {
					self.writeln("  movsd xmm0, qword ptr [rax]")
					self.writeln("  movsd xmm1, qword ptr [rax + 8]")
				}
			}
		} else {
			self.writeln("  mov rdi, qword ptr [rbp - 8]")
			for i := uint(0); i < t.GetByte(); i++ {
				self.writeln("  mov sil, byte ptr [rax + %d]", i)
				self.writeln("  mov byte ptr [rdi + %d], sil", i)
			}
		}
	}
}

// 函数返回值出栈
func (self *CodeGenerator) popRet(t ir.Type) {
	if !ir.IsBasicType(t) {
		if t.GetByte() <= 16 {
			self.writeln("  sub rsp, 16")
			f, s := self.analyseTypeRegister(0, t)
			if t.GetByte() <= 4 {
				if f {
					self.writeln("  mov dword ptr [rsp], eax")
				} else {
					self.writeln("  movss dword ptr [rsp], xmm0")
				}
			} else if t.GetByte() <= 8 {
				if f {
					self.writeln("  mov qword ptr [rsp], rax")
				} else {
					self.writeln("  movsd qword ptr [rsp], xmm0")
				}
			} else {
				if f && s {
					self.writeln("  mov qword ptr [rsp], rax")
					self.writeln("  mov qword ptr [rsp + 8], rdx")
				} else if f {
					self.writeln("  mov qword ptr [rsp], rax")
					self.writeln("  movsd qword ptr [rsp + 8], xmm0")
				} else if s {
					self.writeln("  movsd qword ptr [rsp], xmm0")
					self.writeln("  mov qword ptr [rsp + 8], rax")
				} else {
					self.writeln("  movsd qword ptr [rsp], xmm0")
					self.writeln("  movsd qword ptr [rsp + 8], xmm0")
				}
			}
		}
		self.writeln("  lea rax, [rsp]")
	}
}

// 分析类型前八个字节和后八个字节是否不是全浮点
func (self *CodeGenerator) analyseTypeRegister(offset uint, t ir.Type) (f, e bool) {
	switch typ := t.(type) {
	case ir.IntType, *ir.TypePtr:
		if offset < 8 {
			return true, false
		} else {
			return false, true
		}
	case *ir.TypeFloat:
		return false, false
	case *ir.TypeArray:
		for i := uint(0); offset < 16 && i < typ.Size; i++ {
			offset += typ.GetOffset(i)
			sf, se := self.analyseTypeRegister(offset, typ.Elem)
			f = f || sf
			e = e || se
		}
		return f, e
	case *ir.TypeStruct:
		for i := uint(0); offset < 16 && int(i) < len(typ.Elems); i++ {
			offset += typ.GetOffset(i)
			sf, se := self.analyseTypeRegister(offset, typ.Elems[i])
			f = f || sf
			e = e || se
		}
		return f, e
	default:
		panic("")
	}
}
