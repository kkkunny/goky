package analyse

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/parse"
	"github.com/kkkunny/klang/src/compiler/internal/utils"
	"github.com/kkkunny/stl/table"
	"strings"
)

type Type interface {
	fmt.Stringer
	Equal(Type) bool
}

var (
	None = &typeBase{Name: "none"}

	I8    = &typeBase{Name: "i8"}
	I16   = &typeBase{Name: "i16"}
	I32   = &typeBase{Name: "i32"}
	I64   = &typeBase{Name: "i64"}
	Isize = &typeBase{Name: "isize"}

	U8    = &typeBase{Name: "u8"}
	U16   = &typeBase{Name: "u16"}
	U32   = &typeBase{Name: "u32"}
	U64   = &typeBase{Name: "u64"}
	Usize = &typeBase{Name: "usize"}

	F32 = &typeBase{Name: "f32"}
	F64 = &typeBase{Name: "f64"}

	Bool = &typeBase{Name: "bool"}
)

// typeBase 空类型
type typeBase struct {
	Name string
}

// IsNoneType 是否是空类型
func IsNoneType(t Type) bool {
	return t == None
}

// IsNumberType 是否是数字类型
func IsNumberType(t Type) bool {
	return IsIntType(t) || IsFloatType(t)
}

// IsIntType 是否是整型
func IsIntType(t Type) bool {
	return IsSintType(t) || IsUintType(t)
}

// IsSintType 是否是有符号整型
func IsSintType(t Type) bool {
	return t == I8 || t == I16 || t == I32 || t == I64 || t == Isize
}

// IsUintType 是否是无符号整型
func IsUintType(t Type) bool {
	return t == U8 || t == U16 || t == U32 || t == U64 || t == Usize
}

// IsFloatType 是否是浮点型
func IsFloatType(t Type) bool {
	return t == F32 || t == F64
}

// IsBoolType 是否是布尔类型
func IsBoolType(t Type) bool {
	return t == Bool
}

func (self typeBase) String() string {
	return self.Name
}

func (self typeBase) Equal(t Type) bool {
	if b, ok := t.(*typeBase); ok {
		return self.Name == b.Name
	}
	return false
}

// TypeFunc 函数类型
type TypeFunc struct {
	Ret    Type
	Params []Type
}

// NewFuncType 新建函数类型
func NewFuncType(ret Type, params ...Type) *TypeFunc {
	return &TypeFunc{
		Ret:    ret,
		Params: params,
	}
}

// IsFuncType 是否是函数类型
func IsFuncType(t Type) bool {
	_, ok := t.(*TypeFunc)
	return ok
}

func (self TypeFunc) String() string {
	var buf strings.Builder
	buf.WriteString("func(")
	buf.WriteByte(')')
	buf.WriteString(self.Ret.String())
	return buf.String()
}

func (self TypeFunc) Equal(t Type) bool {
	if f, ok := t.(*TypeFunc); ok {
		if !self.Ret.Equal(f.Ret) || len(self.Params) != len(f.Params) {
			return false
		}
		for i, p := range self.Params {
			if !p.Equal(f.Params[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// TypeArray 数组类型
type TypeArray struct {
	Size uint
	Elem Type
}

// NewArrayType 新建数组类型
func NewArrayType(size uint, elem Type) *TypeArray {
	return &TypeArray{
		Size: size,
		Elem: elem,
	}
}

// IsArrayType 是否是数组类型
func IsArrayType(t Type) bool {
	_, ok := t.(*TypeArray)
	return ok
}

func (self TypeArray) String() string {
	return fmt.Sprintf("[%d]%s", self.Size, self.Elem)
}

func (self TypeArray) Equal(t Type) bool {
	if a, ok := t.(*TypeArray); ok {
		return self.Size == a.Size && self.Elem.Equal(a.Elem)
	}
	return false
}

// TypeTuple 元组类型
type TypeTuple struct {
	Elems []Type
}

// NewTupleType 新建元组类型
func NewTupleType(elems ...Type) *TypeTuple {
	return &TypeTuple{Elems: elems}
}

// IsTupleType 是否是元组类型
func IsTupleType(t Type) bool {
	_, ok := t.(*TypeTuple)
	return ok
}

func (self TypeTuple) String() string {
	types := make([]string, len(self.Elems))
	for i, t := range self.Elems {
		types[i] = t.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(types, ","))
}

func (self TypeTuple) Equal(t Type) bool {
	if t, ok := t.(*TypeTuple); ok {
		if len(self.Elems) != len(t.Elems) {
			return false
		}
		for i, e := range self.Elems {
			if !e.Equal(t.Elems[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// TypeStruct 结构体类型
type TypeStruct struct {
	Fields *table.LinkedHashMap[string, Type]
}

// NewStructType 新建结构体类型
func NewStructType(fields *table.LinkedHashMap[string, Type]) *TypeStruct {
	return &TypeStruct{Fields: fields}
}

// IsStructType 是否是结构体类型
func IsStructType(t Type) bool {
	_, ok := t.(*TypeStruct)
	return ok
}

func (self TypeStruct) String() string {
	var buf strings.Builder
	for iter := self.Fields.Begin(); iter.HasValue(); iter.Next() {
		buf.WriteString(fmt.Sprintf("%s: %s", iter.Key(), iter.Value()))
		if iter.HasNext() {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func (self TypeStruct) Equal(t Type) bool {
	if s, ok := t.(*TypeStruct); ok {
		if self.Fields.Length() != s.Fields.Length() {
			return false
		}
		for iter := self.Fields.Begin(); iter.HasValue(); iter.Next() {
			sk, sv := s.Fields.GetByIndex(iter.Index())
			if iter.Key() != sk || !iter.Value().Equal(sv) {
				return false
			}
		}
		return true
	}
	return false
}

// TypePtr 指针类型
type TypePtr struct {
	Elem Type
}

// NewPtrType 新建指针类型
func NewPtrType(elem Type) *TypePtr {
	return &TypePtr{
		Elem: elem,
	}
}

// IsPtrType 是否是指针类型
func IsPtrType(t Type) bool {
	_, ok := t.(*TypePtr)
	return ok
}

func (self TypePtr) String() string {
	return "*" + self.Elem.String()
}

func (self TypePtr) Equal(t Type) bool {
	if a, ok := t.(*TypePtr); ok {
		return self.Elem.Equal(a.Elem)
	}
	return false
}

// *********************************************************************************************************************

// 类型
func analyseType(ctx *packageContext, ast *parse.Type) (Type, utils.Error) {
	if ast == nil {
		return None, nil
	}
	switch {
	case ast.Ident != nil:
		switch ast.Ident.Name.Value {
		case "i8":
			return I8, nil
		case "i16":
			return I16, nil
		case "i32":
			return I32, nil
		case "i64":
			return I64, nil
		case "isize":
			return Isize, nil
		case "u8":
			return U8, nil
		case "u16":
			return U16, nil
		case "u32":
			return U32, nil
		case "u64":
			return U64, nil
		case "usize":
			return Usize, nil
		case "f32":
			return F32, nil
		case "f64":
			return F64, nil
		case "bool":
			return Bool, nil
		default:
			return nil, utils.Errorf(ast.Ident.Position, "unknown identifier")
		}
	case ast.Func != nil:
		ret, err := analyseType(ctx, ast.Func.Ret)
		if err != nil {
			return nil, err
		}
		params := make([]Type, len(ast.Func.Params.Types))
		var errors []utils.Error
		for i, p := range ast.Func.Params.Types {
			param, err := analyseType(ctx, &p)
			if err != nil {
				errors = append(errors, err)
			} else {
				params[i] = param
			}
		}
		if len(errors) == 0 {
			return NewFuncType(ret, params...), nil
		} else if len(errors) == 1 {
			return nil, errors[0]
		} else {
			return nil, utils.NewMultiError(errors...)
		}
	case ast.Array != nil:
		elem, err := analyseType(ctx, &ast.Array.Elem)
		if err != nil {
			return nil, err
		}
		return NewArrayType(ast.Array.Size, elem), nil
	case ast.Tuple != nil:
		elems := make([]Type, len(ast.Tuple.Types.Types))
		var errors []utils.Error
		for i, e := range ast.Tuple.Types.Types {
			elem, err := analyseType(ctx, &e)
			if err != nil {
				errors = append(errors, err)
			} else {
				elems[i] = elem
			}
		}
		if len(errors) == 0 {
			return NewTupleType(elems...), nil
		} else if len(errors) == 1 {
			return nil, errors[0]
		} else {
			return nil, utils.NewMultiError(errors...)
		}
	case ast.Struct != nil:
		fields := table.NewLinkedHashMap[string, Type]()
		var errors []utils.Error
		for _, f := range ast.Struct.Fields {
			ft, err := analyseType(ctx, &f.Type)
			if err != nil {
				errors = append(errors, err)
			} else if fields.ContainKey(f.Name.Value) {
				errors = append(errors, utils.Errorf(f.Name.Position, "duplicate identifier"))
			} else {
				fields.Set(f.Name.Value, ft)
			}
		}
		if len(errors) == 0 {
			return NewStructType(fields), nil
		} else if len(errors) == 1 {
			return nil, errors[0]
		} else {
			return nil, utils.NewMultiError(errors...)
		}
	case ast.Pointer != nil:
		elem, err := analyseType(ctx, ast.Pointer)
		if err != nil {
			return nil, err
		}
		return NewPtrType(elem), nil
	default:
		panic("")
	}
}
