package generate_ssa

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
)

// 表达式
func (self *Generator) generateExpr(mean analyse.Expr, getValue bool) ir.Value {
	switch expr := mean.(type) {
	case *analyse.Null, *analyse.Integer, *analyse.Float, *analyse.EmptyStruct, *analyse.EmptyArray, *analyse.EmptyTuple:
		return self.generateConstantExpr(mean)
	case *analyse.Binary:
		switch expr.Opera {
		case "+":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewAdd(l, r)
		case "-":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewSub(l, r)
		case "*":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewMul(l, r)
		case "/":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewDiv(l, r)
		case "%":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewMod(l, r)
		case "&":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewAnd(l, r)
		case "|":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewOr(l, r)
		case "^":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewXor(l, r)
		case "<<":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewShl(l, r)
		case ">>":
			l, r := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewShr(l, r)
		case "&&":
			nb, eb := self.block.Belong.NewBlock(), self.block.Belong.NewBlock()
			self.block.NewCondGoto(self.generateExpr(expr.Left, true), nb, eb)
			pb := self.block

			self.block = nb
			nv := self.generateExpr(expr.Right, true)
			self.block.NewGoto(eb)
			nb = self.block

			self.block = eb
			return self.block.NewPhi([]*ir.Block{pb, nb}, []ir.Value{ir.NewInt(ir.I8, 0), nv})
		case "||":
			nb, eb := self.block.Belong.NewBlock(), self.block.Belong.NewBlock()
			self.block.NewCondGoto(self.generateExpr(expr.Left, true), eb, nb)
			pb := self.block

			self.block = nb
			nv := self.generateExpr(expr.Right, true)
			self.block.NewGoto(eb)
			nb = self.block

			self.block = eb
			return self.block.NewPhi([]*ir.Block{pb, nb}, []ir.Value{ir.NewInt(ir.I8, 1), nv})
		default:
			panic("")
		}
	case *analyse.Variable:
		v := self.vars[expr]
		if getValue {
			v = self.block.NewLoad(v)
		}
		return v
	case *analyse.Function:
		return self.vars[expr]
	case *analyse.Call:
		f := self.generateExpr(expr.Func, true)
		args := make([]ir.Value, len(expr.Args))
		for i, a := range expr.Args {
			args[i] = self.generateExpr(a, true)
		}
		if expr.Exit {
			self.doneBeforeFuncEnd()
		}
		call := self.block.NewCall(f, args...)
		if expr.NoReturn {
			self.block.NewUnreachable()
		}
		return call
	case *analyse.Param:
		v := self.vars[expr]
		if getValue {
			v = self.block.NewLoad(v)
		}
		return v
	case *analyse.Array:
		tmp := self.block.NewAlloc(generateType(expr.Type))
		for i, e := range expr.Elems {
			index := self.block.NewArrayIndex(tmp, ir.NewInt(ir.Usize, int64(i)))
			elem := self.generateExpr(e, true)
			self.block.NewStore(elem, index)
		}
		return self.block.NewLoad(tmp)
	case *analyse.Assign:
		switch expr.Opera {
		case "=":
			left, right := self.generateExpr(expr.Left, false), self.generateExpr(expr.Right, true)
			self.block.NewStore(right, left)
			return nil
		default:
			return self.generateExpr(&analyse.Assign{
				Opera: "=",
				Left:  expr.Left,
				Right: &analyse.Binary{
					Opera: expr.Opera[:len(expr.Opera)-1],
					Left:  expr.Left,
					Right: expr.Right,
				},
			}, true)
		}
	case *analyse.Boolean:
		var v int64
		if expr.Value {
			v = 1
		} else {
			v = 0
		}
		return ir.NewInt(ir.I8, v)
	case *analyse.Equal:
		switch expr.Opera {
		case "==":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.equal(left, right)
		case "!=":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			value := self.equal(left, right)
			return self.block.NewXor(value, ir.NewInt(value.GetType().(ir.IntType), 1))
		case "<":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewLt(left, right)
		case "<=":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewLe(left, right)
		case ">":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewGt(left, right)
		case ">=":
			left, right := self.generateExpr(expr.Left, true), self.generateExpr(expr.Right, true)
			return self.block.NewGe(left, right)
		default:
			panic("")
		}
	case *analyse.Unary:
		switch expr.Opera {
		case "!":
			value := self.generateExpr(expr.Value, true)
			return self.block.NewXor(value, ir.NewInt(value.GetType().(ir.IntType), 1))
		case "&":
			return self.generateExpr(expr.Value, false)
		case "*":
			value := self.generateExpr(expr.Value, true)
			if getValue {
				value = self.block.NewLoad(value)
			}
			return value
		default:
			panic("")
		}
	case *analyse.Index:
		fromType := expr.From.GetType()
		switch {
		case analyse.IsArrayType(fromType):
			from, index := self.generateExpr(expr.From, false), self.generateExpr(expr.Index, true)
			var v ir.Value = self.block.NewArrayIndex(from, index)
			if ir.IsTypePtr(from.GetType()) && getValue {
				v = self.block.NewLoad(v)
			}
			return v
		case analyse.IsPtrType(fromType):
			from, index := self.generateExpr(expr.From, true), self.generateExpr(expr.Index, true)
			var v ir.Value = self.block.NewPtrIndex(from, index)
			if ir.IsTypePtr(from.GetType()) && getValue {
				v = self.block.NewLoad(v)
			}
			return v
		case analyse.IsTupleType(fromType):
			from := self.generateExpr(expr.From, false)
			var v ir.Value = self.block.NewStructIndex(from, uint(expr.Index.(*analyse.Integer).Value))
			if ir.IsTypePtr(from.GetType()) && getValue {
				v = self.block.NewLoad(v)
			}
			return v
		default:
			panic("")
		}
	case *analyse.Select:
		cond := self.generateExpr(expr.Cond, true)
		tb, fb, eb := self.block.Belong.NewBlock(), self.block.Belong.NewBlock(), self.block.Belong.NewBlock()
		self.block.NewCondGoto(cond, tb, fb)

		self.block = tb
		tv := self.generateExpr(expr.True, getValue)
		self.block.NewGoto(eb)

		self.block = fb
		fv := self.generateExpr(expr.False, getValue)
		self.block.NewGoto(eb)

		self.block = eb
		return self.block.NewPhi([]*ir.Block{tb, fb}, []ir.Value{tv, fv})
	case *analyse.Tuple:
		tmp := self.block.NewAlloc(generateType(expr.Type))
		for i, e := range expr.Elems {
			index := self.block.NewStructIndex(tmp, uint(i))
			elem := self.generateExpr(e, true)
			self.block.NewStore(elem, index)
		}
		return self.block.NewLoad(tmp)
	case *analyse.Struct:
		tmp := self.block.NewAlloc(generateType(expr.Type))
		for i, e := range expr.Fields {
			index := self.block.NewStructIndex(tmp, uint(i))
			elem := self.generateExpr(e, true)
			self.block.NewStore(elem, index)
		}
		return self.block.NewLoad(tmp)
	case *analyse.GetField:
		f := self.generateExpr(expr.From, false)
		var index uint
		for iter := expr.From.GetType().(*analyse.TypeStruct).Fields.Begin(); iter.HasValue(); iter.Next() {
			if iter.Key() == expr.Index {
				break
			}
			index++
		}
		var v ir.Value = self.block.NewStructIndex(f, index)
		if ir.IsTypePtr(f.GetType()) && getValue {
			v = self.block.NewLoad(v)
		}
		return v
	case *analyse.Covert:
		from := self.generateExpr(expr.From, true)
		meanFt, meanTo := expr.From.GetType(), expr.To
		to := generateType(expr.GetType())
		switch {
		case analyse.IsIntType(meanFt) && analyse.IsIntType(meanTo):
			return self.block.NewItoi(from, to.(ir.IntType))
		case analyse.IsIntType(meanFt) && analyse.IsFloatType(meanTo):
			return self.block.NewItof(from, to.(*ir.TypeFloat))
		case analyse.IsFloatType(meanFt) && analyse.IsIntType(meanTo):
			return self.block.NewFtoi(from, to.(ir.IntType))
		case analyse.IsIntType(meanFt) && (analyse.IsPtrType(meanTo) || analyse.IsFuncType(meanTo)):
			return self.block.NewItop(from, to.(*ir.TypePtr))
		case (analyse.IsPtrType(meanFt) || analyse.IsFuncType(meanFt)) && analyse.IsIntType(meanTo):
			return self.block.NewPtoi(from, to.(ir.IntType))
		case (analyse.IsPtrType(meanFt) || analyse.IsFuncType(meanFt)) && (analyse.IsPtrType(meanTo) || analyse.IsFuncType(meanTo)):
			return self.block.NewPtop(from, to.(*ir.TypePtr))
		default:
			panic("")
		}
	case *analyse.GlobalVariable:
		v := self.vars[expr]
		if getValue {
			v = self.block.NewLoad(v)
		}
		return v
	default:
		panic("")
	}
}

// 常量表达式
func (self *Generator) generateConstantExpr(mean analyse.Expr) ir.Constant {
	switch expr := mean.(type) {
	case *analyse.Null:
		return ir.NewEmpty(generateType(expr.Type))
	case *analyse.Integer:
		return ir.NewInt(generateType(expr.Type).(ir.IntType), expr.Value)
	case *analyse.Float:
		return ir.NewFloat(generateType(expr.Type).(*ir.TypeFloat), expr.Value)
	case *analyse.EmptyArray:
		return ir.NewEmpty(generateType(expr.Type))
	case *analyse.Array:
		elems := make([]ir.Constant, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateConstantExpr(e)
		}
		return ir.NewArray(elems...)
	case *analyse.EmptyTuple:
		return ir.NewEmpty(generateType(expr.Type))
	case *analyse.Tuple:
		elems := make([]ir.Constant, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateConstantExpr(e)
		}
		return ir.NewStruct(elems...)
	case *analyse.EmptyStruct:
		return ir.NewEmpty(generateType(expr.Type))
	case *analyse.Struct:
		elems := make([]ir.Constant, len(expr.Fields))
		for i, e := range expr.Fields {
			elems[i] = self.generateConstantExpr(e)
		}
		return ir.NewStruct(elems...)
	default:
		panic("")
	}
}

// 比较
func (self *Generator) equal(left, right ir.Value) ir.Value {
	if !left.GetType().Equal(right.GetType()) {
		panic("")
	}
	typ := left.GetType()
	switch t := typ.(type) {
	case ir.NumberType, *ir.TypePtr:
		return self.block.NewEq(left, right)
	case *ir.TypeArray:
		if t.Size == 0 {
			return ir.NewInt(ir.I8, 1)
		}
		i := self.block.NewAlloc(ir.Usize)
		self.block.NewStore(ir.NewInt(ir.Usize, 0), i)
		cb := self.block.Belong.NewBlock()
		self.block.NewGoto(cb)

		self.block = cb
		iv := self.block.NewLoad(i)
		lb, eb := self.block.Belong.NewBlock(), self.block.Belong.NewBlock()
		self.block.NewCondGoto(self.block.NewLt(iv, ir.NewInt(ir.Usize, int64(t.Size))), lb, eb)

		self.block = lb
		l, r := self.block.NewArrayIndex(left, iv), self.block.NewArrayIndex(right, iv)
		self.block.NewStore(self.block.NewAdd(iv, ir.NewInt(ir.Usize, 1)), i)
		self.block.NewCondGoto(self.equal(l, r), cb, eb)
		lb = self.block

		self.block = eb
		return self.block.NewPhi([]*ir.Block{cb, lb}, []ir.Value{ir.NewInt(ir.I8, 1), ir.NewInt(ir.I8, 0)})
	case *ir.TypeStruct:
		if len(t.Elems) == 0 {
			return ir.NewInt(ir.I8, 1)
		}
		blocks := make([]*ir.Block, len(t.Elems))
		values := make([]ir.Value, len(t.Elems))
		endBlock := self.block.Belong.NewBlock()
		for i := range t.Elems {
			l, r := self.block.NewStructIndex(left, uint(i)), self.block.NewStructIndex(right, uint(i))
			v := self.equal(l, r)
			blocks[i], values[i] = self.block, v
			if i < len(t.Elems)-1 {
				nextBlock := self.block.Belong.NewBlock()
				self.block.NewCondGoto(v, nextBlock, endBlock)
				self.block = nextBlock
			} else {
				self.block.NewGoto(endBlock)
				self.block = endBlock
			}
		}
		return self.block.NewPhi(blocks, values)
	default:
		panic("")
	}
}
