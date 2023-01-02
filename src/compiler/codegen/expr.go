package codegen

import (
	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/klang/src/compiler/analyse"
	stlutil "github.com/kkkunny/stl/util"
	"unsafe"
)

// 表达式
func (self *CodeGenerator) codegenExpr(mean analyse.Expr, getValue bool) llvm.Value {
	switch expr := mean.(type) {
	case *analyse.Null, *analyse.Integer, *analyse.Float, *analyse.Boolean, *analyse.EmptyStruct, *analyse.EmptyArray, *analyse.EmptyTuple:
		return self.codegenConstantExpr(mean)
	case *analyse.Binary:
		switch expr.Opera {
		case "+":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNSWAdd(l, r, "")
			} else if analyse.IsUintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNUWAdd(l, r, "")
			} else {
				return self.builder.CreateFAdd(l, r, "")
			}
		case "-":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNSWSub(l, r, "")
			} else if analyse.IsUintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNUWSub(l, r, "")
			} else {
				return self.builder.CreateFSub(l, r, "")
			}
		case "*":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNSWMul(l, r, "")
			} else if analyse.IsUintTypeAndSon(expr.GetType()) {
				return self.builder.CreateNUWMul(l, r, "")
			} else {
				return self.builder.CreateFMul(l, r, "")
			}
		case "/":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateSDiv(l, r, "")
			} else if analyse.IsUintTypeAndSon(expr.GetType()) {
				return self.builder.CreateUDiv(l, r, "")
			} else {
				return self.builder.CreateFDiv(l, r, "")
			}
		case "%":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateSRem(l, r, "")
			} else if analyse.IsUintTypeAndSon(expr.GetType()) {
				return self.builder.CreateURem(l, r, "")
			} else {
				return self.builder.CreateFRem(l, r, "")
			}
		case "&":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			return self.builder.CreateAnd(l, r, "")
		case "|":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			return self.builder.CreateOr(l, r, "")
		case "^":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			return self.builder.CreateXor(l, r, "")
		case "<<":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			return self.builder.CreateShl(l, r, "")
		case ">>":
			l, r := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.GetType()) {
				return self.builder.CreateAShr(l, r, "")
			} else {
				return self.builder.CreateLShr(l, r, "")
			}
		case "&&":
			nb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
			self.builder.CreateCondBr(self.codegenExpr(expr.Left, true), nb, eb)
			pb := self.builder.GetInsertBlock()

			self.builder.SetInsertPointAtEnd(nb)
			nv := self.codegenExpr(expr.Right, true)
			self.builder.CreateBr(eb)
			nb = self.builder.GetInsertBlock()

			self.builder.SetInsertPointAtEnd(eb)
			phi := self.builder.CreatePHI(nv.Type(), "")
			phi.AddIncoming([]llvm.Value{v_false, nv}, []llvm.BasicBlock{pb, nb})
			return phi
		case "||":
			nb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
			self.builder.CreateCondBr(self.codegenExpr(expr.Left, true), eb, nb)
			pb := self.builder.GetInsertBlock()

			self.builder.SetInsertPointAtEnd(nb)
			nv := self.codegenExpr(expr.Right, true)
			self.builder.CreateBr(eb)
			nb = self.builder.GetInsertBlock()

			self.builder.SetInsertPointAtEnd(eb)
			phi := self.builder.CreatePHI(nv.Type(), "")
			phi.AddIncoming([]llvm.Value{v_true, nv}, []llvm.BasicBlock{pb, nb})
			return phi
		default:
			panic("")
		}
	case *analyse.Variable:
		v := self.vars[expr]
		if getValue {
			v = self.builder.CreateLoad(v.Type().ElementType(), v, "")
		}
		return v
	case *analyse.Function:
		return self.vars[expr]
	case *analyse.Method:
		return self.vars[expr.Func]
	case *analyse.FuncCall:
		f := self.codegenExpr(expr.Func, true)
		args := make([]llvm.Value, len(expr.Args))
		for i, a := range expr.Args {
			args[i] = self.codegenExpr(a, true)
		}
		if meanFunc, ok := expr.Func.(*analyse.Function); ok {
			if meanFunc.Exit {
				self.doneBeforeFuncEnd()
			}
		}
		call := self.builder.CreateCall(f.Type().ReturnType(), f, args, "")
		if meanFunc, ok := expr.Func.(*analyse.Function); ok {
			if meanFunc.NoReturn || meanFunc.Exit {
				self.doneBeforeFuncEnd()
			}
		}
		return call
	case *analyse.MethodCall:
		f := self.codegenExpr(expr.Method.Func, true)
		args := make([]llvm.Value, len(expr.Args)+1)
		if analyse.IsPtrType(expr.Method.Self.GetType()) {
			args[0] = self.codegenExpr(expr.Method.Self, true)
		} else if expr.Method.Self.GetMut() {
			args[0] = self.codegenExpr(expr.Method.Self, false)
		} else {
			selfArg := self.codegenExpr(expr.Method.Self, true)
			args[0] = self.builder.CreateAlloca(selfArg.Type(), "")
			self.builder.CreateStore(selfArg, args[0])
		}
		for i, a := range expr.Args {
			args[i+1] = self.codegenExpr(a, true)
		}
		if expr.Method.Func.Exit {
			self.doneBeforeFuncEnd()
		}
		call := self.builder.CreateCall(f.Type().ReturnType(), f, args, "")
		if expr.Method.Func.NoReturn || expr.Method.Func.Exit {
			self.builder.CreateUnreachable()
		}
		return call
	case *analyse.Param:
		v := self.vars[expr]
		if getValue {
			v = self.builder.CreateLoad(v.Type().ElementType(), v, "")
		}
		return v
	case *analyse.Assign:
		switch expr.Opera {
		case "=":
			left, right := self.codegenExpr(expr.Left, false), self.codegenExpr(expr.Right, true)
			self.builder.CreateStore(right, left)
			return llvm.Value{}
		default:
			return self.codegenExpr(&analyse.Assign{
				Opera: "=",
				Left:  expr.Left,
				Right: &analyse.Binary{
					Opera: expr.Opera[:len(expr.Opera)-1],
					Left:  expr.Left,
					Right: expr.Right,
				},
			}, true)
		}
	case *analyse.Equal:
		switch expr.Opera {
		case "==":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			return self.equal(left, right)
		case "!=":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			left = self.equal(left, right)
			return self.builder.CreateXor(left, llvm.ConstInt(left.Type(), 1, true), "")
		case "<":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntSLT, left, right, "")
			} else if analyse.IsUintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntULT, left, right, "")
			} else {
				return self.builder.CreateFCmp(llvm.FloatOLT, left, right, "")
			}
		case "<=":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntSLE, left, right, "")
			} else if analyse.IsUintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntULE, left, right, "")
			} else {
				return self.builder.CreateFCmp(llvm.FloatOLE, left, right, "")
			}
		case ">":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntSGT, left, right, "")
			} else if analyse.IsUintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntUGT, left, right, "")
			} else {
				return self.builder.CreateFCmp(llvm.FloatOGT, left, right, "")
			}
		case ">=":
			left, right := self.codegenExpr(expr.Left, true), self.codegenExpr(expr.Right, true)
			if analyse.IsSintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntSGE, left, right, "")
			} else if analyse.IsUintTypeAndSon(expr.Left.GetType()) {
				return self.builder.CreateICmp(llvm.IntUGE, left, right, "")
			} else {
				return self.builder.CreateFCmp(llvm.FloatOGE, left, right, "")
			}
		default:
			panic("")
		}
	case *analyse.Unary:
		switch expr.Opera {
		case "!":
			left := self.codegenExpr(expr.Value, true)
			return self.builder.CreateXor(left, llvm.ConstInt(left.Type(), 1, true), "")
		case "&":
			return self.codegenExpr(expr.Value, false)
		case "*":
			value := self.codegenExpr(expr.Value, true)
			if getValue {
				value = self.builder.CreateLoad(value.Type().ElementType(), value, "")
			}
			return value
		default:
			panic("")
		}
	case *analyse.Index:
		fromType := expr.From.GetType()
		switch {
		case analyse.IsArrayTypeAndSon(fromType):
			from, index := self.codegenExpr(expr.From, false), self.codegenExpr(expr.Index, true)
			return self.createArrayIndex(from, index, getValue)
		case analyse.IsPtrTypeAndSon(fromType):
			from, index := self.codegenExpr(expr.From, true), self.codegenExpr(expr.Index, true)
			return self.createPointerIndex(from, index, getValue)
		case analyse.IsTupleTypeAndSon(fromType):
			from := self.codegenExpr(expr.From, false)
			index := expr.Index.(*analyse.Integer).Value
			return self.createStructIndex(from, uint(index), getValue)
		default:
			panic("")
		}
	case *analyse.Select:
		cond := self.codegenExpr(expr.Cond, true)
		tb, fb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
		self.builder.CreateCondBr(cond, tb, fb)

		self.builder.SetInsertPointAtEnd(tb)
		tv := self.codegenExpr(expr.True, getValue)
		self.builder.CreateBr(eb)

		self.builder.SetInsertPointAtEnd(fb)
		fv := self.codegenExpr(expr.False, getValue)
		self.builder.CreateBr(eb)

		self.builder.SetInsertPointAtEnd(eb)
		phi := self.builder.CreatePHI(tv.Type(), "")
		phi.AddIncoming([]llvm.Value{tv, fv}, []llvm.BasicBlock{tb, fb})
		return phi
	case *analyse.Array:
		isConst := true
		elems := make([]llvm.Value, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.codegenExpr(e, true)
			if !elems[i].IsConstant() {
				isConst = false
			}
		}
		if isConst {
			return llvm.ConstArray(self.codegenType(expr.Type), elems)
		} else {
			tmp := self.builder.CreateAlloca(self.codegenType(expr.Type), "")
			for i, e := range elems {
				index := self.createArrayIndex(tmp, llvm.ConstInt(t_size, uint64(i), false), false)
				self.builder.CreateStore(e, index)
			}
			return self.builder.CreateLoad(tmp.Type().ElementType(), tmp, "")
		}
	case *analyse.Tuple:
		isConst := true
		elems := make([]llvm.Value, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.codegenExpr(e, true)
			if !elems[i].IsConstant() {
				isConst = false
			}
		}
		if isConst {
			return self.ctx.ConstStruct(elems, false)
		} else {
			tmp := self.builder.CreateAlloca(self.codegenType(expr.Type), "")
			for i, e := range elems {
				index := self.createStructIndex(tmp, uint(i), false)
				self.builder.CreateStore(e, index)
			}
			return self.builder.CreateLoad(tmp.Type().ElementType(), tmp, "")
		}
	case *analyse.Struct:
		isConst := true
		elems := make([]llvm.Value, len(expr.Fields))
		for i, e := range expr.Fields {
			elems[i] = self.codegenExpr(e, true)
			if !elems[i].IsConstant() {
				isConst = false
			}
		}
		if isConst {
			return self.ctx.ConstStruct(elems, false)
		} else {
			tmp := self.builder.CreateAlloca(self.codegenType(expr.Type), "")
			for i, e := range elems {
				index := self.createStructIndex(tmp, uint(i), false)
				self.builder.CreateStore(e, index)
			}
			return self.builder.CreateLoad(tmp.Type().ElementType(), tmp, "")
		}
	case *analyse.GetField:
		f := self.codegenExpr(expr.From, false)
		var index uint
		for iter := analyse.GetBaseType(expr.From.GetType()).(*analyse.TypeStruct).Fields.Begin(); iter.HasValue(); iter.Next() {
			if iter.Key() == expr.Index {
				break
			}
			index++
		}
		return self.createStructIndex(f, index, getValue)
	case *analyse.Covert:
		from := self.codegenExpr(expr.From, true)
		meanFt, meanTo := expr.From.GetType(), expr.To
		to := self.codegenType(expr.GetType())
		switch {
		case analyse.GetDepthBaseType(meanFt).Equal(analyse.GetDepthBaseType(meanTo)):
			return from
		case analyse.IsIntTypeAndSon(meanFt) && analyse.IsIntTypeAndSon(meanTo):
			return self.builder.CreateIntCast(from, to, "")
		case analyse.IsFloatTypeAndSon(meanFt) && analyse.IsFloatTypeAndSon(meanTo):
			return self.builder.CreateFPCast(from, to, "")
		case analyse.IsSintTypeAndSon(meanFt) && analyse.IsFloatTypeAndSon(meanTo):
			return self.builder.CreateSIToFP(from, to, "")
		case analyse.IsUintTypeAndSon(meanFt) && analyse.IsFloatTypeAndSon(meanTo):
			return self.builder.CreateUIToFP(from, to, "")
		case analyse.IsFloatTypeAndSon(meanFt) && analyse.IsSintTypeAndSon(meanTo):
			return self.builder.CreateFPToSI(from, to, "")
		case analyse.IsFloatTypeAndSon(meanFt) && analyse.IsUintTypeAndSon(meanTo):
			return self.builder.CreateFPToUI(from, to, "")
		case analyse.GetBaseType(meanFt).Equal(analyse.Usize) && (analyse.IsPtrTypeAndSon(meanTo) || analyse.IsFuncTypeAndSon(meanTo)):
			return self.builder.CreateIntToPtr(from, to, "")
		case (analyse.IsPtrTypeAndSon(meanFt) || analyse.IsFuncTypeAndSon(meanFt)) && analyse.GetBaseType(meanTo).Equal(analyse.Usize):
			return self.builder.CreatePtrToInt(from, to, "")
		case (analyse.IsPtrTypeAndSon(meanFt) || analyse.IsFuncTypeAndSon(meanFt)) && (analyse.IsPtrTypeAndSon(meanTo) || analyse.IsFuncTypeAndSon(meanTo)):
			return self.builder.CreatePointerCast(from, to, "")
		default:
			panic("")
		}
	case *analyse.GlobalVariable:
		v := self.vars[expr]
		if getValue {
			v = self.builder.CreateLoad(v.Type().ElementType(), v, "")
		}
		return v
	case *analyse.GetTypeBytes:
		return llvm.SizeOf(self.codegenType(expr.Type))
	default:
		panic("")
	}
}

// 常量表达式
func (self *CodeGenerator) codegenConstantExpr(mean analyse.Expr) llvm.Value {
	switch expr := mean.(type) {
	case *analyse.Null:
		return llvm.ConstPointerNull(self.codegenType(expr.Type))
	case *analyse.Integer:
		value := *(*uint64)(unsafe.Pointer(&expr.Value))
		return llvm.ConstInt(self.codegenType(expr.Type), value, analyse.IsSintType(expr.Type))
	case *analyse.Float:
		return llvm.ConstFloat(self.codegenType(expr.Type), expr.Value)
	case *analyse.Boolean:
		return stlutil.Ternary(expr.Value, v_true, v_false)
	case *analyse.EmptyArray, *analyse.EmptyTuple, *analyse.EmptyStruct:
		return llvm.ConstNull(self.codegenType(expr.GetType()))
	case *analyse.Array:
		elems := make([]llvm.Value, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.codegenConstantExpr(e)
		}
		return llvm.ConstArray(self.codegenType(expr.Type), elems)
	case *analyse.Tuple:
		elems := make([]llvm.Value, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.codegenConstantExpr(e)
		}
		return self.ctx.ConstStruct(elems, false)
	case *analyse.Struct:
		elems := make([]llvm.Value, len(expr.Fields))
		for i, e := range expr.Fields {
			elems[i] = self.codegenConstantExpr(e)
		}
		return self.ctx.ConstStruct(elems, false)
	default:
		panic("")
	}
}

// 比较
func (self *CodeGenerator) equal(left, right llvm.Value) llvm.Value {
	switch left.Type().TypeKind() {
	case llvm.IntegerTypeKind, llvm.PointerTypeKind, llvm.FunctionTypeKind:
		return self.builder.CreateICmp(llvm.IntEQ, left, right, "")
	case llvm.FloatTypeKind:
		return self.builder.CreateFCmp(llvm.FloatOEQ, left, right, "")
	case llvm.ArrayTypeKind:
		if left.Type().ArrayLength() == 0 {
			return llvm.ConstInt(self.ctx.Int8Type(), 1, true)
		}
		i := self.builder.CreateAlloca(self.codegenType(analyse.Usize), "")
		self.builder.CreateStore(llvm.ConstInt(i.Type().ElementType(), 0, false), i)
		cb := llvm.AddBasicBlock(self.function, "")
		self.builder.CreateBr(cb)

		self.builder.SetInsertPointAtEnd(cb)
		iv := self.builder.CreateLoad(i.Type().ElementType(), i, "")
		lb, eb := llvm.AddBasicBlock(self.function, ""), llvm.AddBasicBlock(self.function, "")
		lt := self.builder.CreateICmp(llvm.IntULT, iv, llvm.ConstInt(iv.Type(), uint64(left.Type().ArrayLength()), false), "")
		self.builder.CreateCondBr(lt, lb, eb)

		self.builder.SetInsertPointAtEnd(lb)
		l, r := self.createArrayIndex(left, iv, true), self.createArrayIndex(right, iv, true)
		self.builder.CreateStore(self.builder.CreateNUWAdd(iv, llvm.ConstInt(iv.Type(), 1, false), ""), i)
		self.builder.CreateCondBr(self.equal(l, r), cb, eb)
		lb = self.builder.GetInsertBlock()

		self.builder.SetInsertPointAtEnd(eb)
		phi := self.builder.CreatePHI(t_bool, "")
		phi.AddIncoming([]llvm.Value{v_true, v_false}, []llvm.BasicBlock{cb, lb})
		return phi
	case llvm.StructTypeKind:
		elemCount := left.Type().StructElementTypesCount()
		if elemCount == 0 {
			return llvm.ConstInt(self.ctx.Int8Type(), 1, true)
		}
		blocks := make([]llvm.BasicBlock, elemCount)
		values := make([]llvm.Value, elemCount)
		eb := llvm.AddBasicBlock(self.function, "")
		for i := range left.Type().StructElementTypes() {
			l, r := self.createStructIndex(left, uint(i), true), self.createStructIndex(right, uint(i), true)
			v := self.equal(l, r)
			blocks[i], values[i] = self.builder.GetInsertBlock(), v
			if i < elemCount-1 {
				nb := llvm.AddBasicBlock(self.function, "")
				self.builder.CreateCondBr(v, nb, eb)
				self.builder.SetInsertPointAtEnd(nb)
			} else {
				self.builder.CreateBr(eb)
				self.builder.SetInsertPointAtEnd(eb)
			}
		}
		phi := self.builder.CreatePHI(t_bool, "")
		phi.AddIncoming(values, blocks)
		return phi
	default:
		panic("")
	}
}
