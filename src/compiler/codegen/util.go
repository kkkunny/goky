package codegen

import (
	"github.com/kkkunny/Sim/src/compiler/utils"
	"github.com/kkkunny/go-llvm"
)

var (
	t_bool llvm.Type
	t_size llvm.Type

	v_true  llvm.Value
	v_false llvm.Value
)

func (self CodeGenerator) init() {
	t_bool = self.ctx.Int8Type()
	t_size = self.ctx.IntType(int(utils.PtrByte * 8))

	v_true = llvm.ConstInt(t_bool, 1, true)
	v_false = llvm.ConstInt(t_bool, 0, true)
}

func (self *CodeGenerator) createArrayIndex(v llvm.Value, i llvm.Value, getValue bool) llvm.Value {
	if v.Type().TypeKind() == llvm.PointerTypeKind {
		value := self.builder.CreateInBoundsGEP(v.Type().ElementType(), v, []llvm.Value{llvm.ConstInt(t_size, 0, false), i}, "")
		if getValue {
			value = self.builder.CreateLoad(value.Type().ElementType(), value, "")
		}
		return value
	} else {
		return self.builder.CreateExtractElement(v, i, "")
	}
}

func (self *CodeGenerator) createPointerIndex(v llvm.Value, i llvm.Value, getValue bool) llvm.Value {
	value := self.builder.CreateInBoundsGEP(v.Type().ElementType(), v, []llvm.Value{i}, "")
	if getValue {
		value = self.builder.CreateLoad(value.Type().ElementType(), value, "")
	}
	return value
}

func (self *CodeGenerator) createStructIndex(v llvm.Value, i uint, getValue bool) llvm.Value {
	if v.Type().TypeKind() == llvm.PointerTypeKind {
		value := self.builder.CreateStructGEP(v.Type().ElementType(), v, int(i), "")
		if getValue {
			value = self.builder.CreateLoad(value.Type().ElementType(), value, "")
		}
		return value
	} else {
		return self.builder.CreateExtractValue(v, int(i), "")
	}
}
