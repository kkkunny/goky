package generate_ssa

import (
	"github.com/kkkunny/klang/src/compiler/analyse"
	"github.com/kkkunny/klang/src/compiler/ir_ssa"
)

// 类型
func generateType(mean analyse.Type) ir.Type {
	switch typ := mean.(type) {
	case *analyse.TypeFunc:
		ret := generateType(typ.Ret)
		params := make([]ir.Type, len(typ.Params))
		for i, p := range typ.Params {
			params[i] = generateType(p)
		}
		return ir.NewPtrType(ir.NewFuncType(ret, params...))
	case *analyse.TypeArray:
		elem := generateType(typ.Elem)
		return ir.NewArrayType(typ.Size, elem)
	case *analyse.TypeTuple:
		elems := make([]ir.Type, len(typ.Elems))
		for i, e := range typ.Elems {
			elems[i] = generateType(e)
		}
		return ir.NewStructType(elems...)
	case *analyse.TypeStruct:
		elems := make([]ir.Type, typ.Fields.Length())
		for iter := typ.Fields.Begin(); iter.HasValue(); iter.Next() {
			elems[iter.Index()] = generateType(iter.Value())
		}
		return ir.NewStructType(elems...)
	case *analyse.TypePtr:
		return ir.NewPtrType(generateType(typ.Elem))
	default:
		switch {
		case analyse.IsNoneType(typ):
			return ir.None
		case analyse.IsSintType(typ):
			switch typ {
			case analyse.I8:
				return ir.I8
			case analyse.I16:
				return ir.I16
			case analyse.I32:
				return ir.I32
			case analyse.I64:
				return ir.I64
			case analyse.Isize:
				return ir.Isize
			default:
				panic("")
			}
		case analyse.IsUintType(typ):
			switch typ {
			case analyse.U8:
				return ir.U8
			case analyse.U16:
				return ir.U16
			case analyse.U32:
				return ir.U32
			case analyse.U64:
				return ir.U64
			case analyse.Usize:
				return ir.Usize
			default:
				panic("")
			}
		case analyse.IsFloatType(typ):
			switch typ {
			case analyse.F32:
				return ir.F32
			case analyse.F64:
				return ir.F64
			default:
				panic("")
			}
		case analyse.IsBoolType(typ):
			return ir.I8
		default:
			panic("")
		}
	}
}
