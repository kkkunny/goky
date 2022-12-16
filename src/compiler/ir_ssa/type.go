package ir

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/utils"
	"math"
	"strings"
)

// Type 类型
type Type interface {
	fmt.Stringer
	Equal(Type) bool
	GetByte() uint
	GetAlign() uint
}

var (
	None = &TypeNone{}

	I8    = &TypeSint{Bytes: 1}
	I16   = &TypeSint{Bytes: 2}
	I32   = &TypeSint{Bytes: 4}
	I64   = &TypeSint{Bytes: 8}
	Isize = &TypeSint{Bytes: utils.PtrByte}

	U8    = &TypeUint{Bytes: 1}
	U16   = &TypeUint{Bytes: 2}
	U32   = &TypeUint{Bytes: 4}
	U64   = &TypeUint{Bytes: 8}
	Usize = &TypeUint{Bytes: utils.PtrByte}

	F32 = &TypeFloat{Bytes: 4}
	F64 = &TypeFloat{Bytes: 8}
)

// BasicType 基础类型
type BasicType interface {
	Type
	basic()
}

// IsBasicType 是否是基础类型
func IsBasicType(t Type) bool {
	_, ok := t.(BasicType)
	return ok
}

// TypeNone 空类型
type TypeNone struct{}

// IsNoneType 是否是空类型
func IsNoneType(t Type) bool {
	_, ok := t.(*TypeNone)
	return ok
}

func (self TypeNone) String() string {
	return "none"
}

func (self TypeNone) Equal(t Type) bool {
	return IsNoneType(t)
}

func (self TypeNone) GetByte() uint {
	panic("")
}

func (self TypeNone) GetAlign() uint {
	panic("")
}

func (self TypeNone) basic() {}

// NumberType 数字类型
type NumberType interface {
	Type
	number()
}

// IsNumberType 是否是数字类型
func IsNumberType(t Type) bool {
	_, ok := t.(NumberType)
	return ok
}

// IntType 整数
type IntType interface {
	NumberType
	int()
}

// IsIntType 是否是整型
func IsIntType(t Type) bool {
	_, ok := t.(IntType)
	return ok
}

// TypeSint 有符号整型
type TypeSint struct {
	Bytes uint
}

// IsSintType 是否是有符号整型
func IsSintType(t Type) bool {
	_, ok := t.(*TypeSint)
	return ok
}

func (self TypeSint) String() string {
	return fmt.Sprintf("i%d", self.Bytes*8)
}

func (self TypeSint) Equal(t Type) bool {
	if s, ok := t.(*TypeSint); ok {
		return self.Bytes == s.Bytes
	}
	return false
}

func (self TypeSint) GetByte() uint {
	return self.Bytes
}

func (self TypeSint) GetAlign() uint {
	return self.Bytes
}

func (self TypeSint) number() {}

func (self TypeSint) int() {}

func (self TypeSint) basic() {}

// TypeUint 无符号整型
type TypeUint struct {
	Bytes uint
}

// IsUintType 是否是无符号整型
func IsUintType(t Type) bool {
	_, ok := t.(*TypeUint)
	return ok
}

func (self TypeUint) String() string {
	return fmt.Sprintf("u%d", self.Bytes*8)
}

func (self TypeUint) Equal(t Type) bool {
	if s, ok := t.(*TypeUint); ok {
		return self.Bytes == s.Bytes
	}
	return false
}

func (self TypeUint) GetByte() uint {
	return self.Bytes
}

func (self TypeUint) GetAlign() uint {
	return self.Bytes
}

func (self TypeUint) number() {}

func (self TypeUint) int() {}

func (self TypeUint) basic() {}

// TypeFloat 浮点型
type TypeFloat struct {
	Bytes uint
}

// IsFloatType 是否是浮点型
func IsFloatType(t Type) bool {
	_, ok := t.(*TypeFloat)
	return ok
}

func (self TypeFloat) String() string {
	return fmt.Sprintf("f%d", self.Bytes*8)
}

func (self TypeFloat) Equal(t Type) bool {
	if s, ok := t.(*TypeFloat); ok {
		return self.Bytes == s.Bytes
	}
	return false
}

func (self TypeFloat) GetByte() uint {
	return self.Bytes
}

func (self TypeFloat) GetAlign() uint {
	return self.Bytes
}

func (self TypeFloat) number() {}

func (self TypeFloat) basic() {}

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
	for i, p := range self.Params {
		buf.WriteString(p.String())
		if i < len(self.Params)-1 {
			buf.WriteString(",")
		}
	}
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

func (self TypeFunc) GetByte() uint {
	return utils.PtrByte
}

func (self TypeFunc) GetAlign() uint {
	return utils.PtrByte
}

func (self TypeFunc) basic() {}

// TypePtr 指针类型
type TypePtr struct {
	Elem Type
}

// NewPtrType 新建指针类型
func NewPtrType(ret Type) *TypePtr {
	return &TypePtr{
		Elem: ret,
	}
}

// IsTypePtr 是否是指针类型
func IsTypePtr(t Type) bool {
	_, ok := t.(*TypePtr)
	return ok
}

func (self TypePtr) String() string {
	return fmt.Sprintf("*%s", self.Elem)
}

func (self TypePtr) Equal(t Type) bool {
	if f, ok := t.(*TypePtr); ok {
		return self.Elem.Equal(f.Elem)
	}
	return false
}

func (self TypePtr) GetByte() uint {
	return utils.PtrByte
}

func (self TypePtr) GetAlign() uint {
	return utils.PtrByte
}

func (self TypePtr) basic() {}

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

func (self TypeArray) GetByte() uint {
	return self.Elem.GetByte() * self.Size
}

func (self TypeArray) GetAlign() uint {
	return self.Elem.GetAlign()
}

func (self TypeArray) GetOffset(i uint) uint {
	if i >= self.Size {
		panic("")
	}
	return self.Elem.GetByte() * i
}

// TypeStruct 结构体类型
type TypeStruct struct {
	Elems []Type
}

// NewStructType 新建结构体类型
func NewStructType(elems ...Type) *TypeStruct {
	return &TypeStruct{Elems: elems}
}

// IsStructType 是否是结构体类型
func IsStructType(t Type) bool {
	_, ok := t.(*TypeStruct)
	return ok
}

func (self TypeStruct) String() string {
	types := make([]string, len(self.Elems))
	for i, t := range self.Elems {
		types[i] = t.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(types, ", "))
}

func (self TypeStruct) Equal(t Type) bool {
	if t, ok := t.(*TypeStruct); ok {
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

func (self TypeStruct) GetByte() (offset uint) {
	align := uint(math.Max(float64(utils.AlignByte), float64(self.GetAlign())))
	for _, e := range self.Elems {
		offset = utils.AlignTo(offset, align)
		offset += e.GetByte()
	}
	return utils.AlignTo(offset, align)
}

func (self TypeStruct) GetOffset(i uint) (offset uint) {
	if i >= uint(len(self.Elems)) {
		panic("")
	}
	align := uint(math.Max(float64(utils.AlignByte), float64(self.GetAlign())))
	for ii, e := range self.Elems {
		offset = utils.AlignTo(offset, align)
		if uint(ii) == i {
			break
		}
		offset += e.GetByte()
	}
	return offset
}

func (self TypeStruct) GetAlign() (align uint) {
	align = 1
	for _, e := range self.Elems {
		align = uint(math.Max(float64(align), float64(e.GetAlign())))
	}
	return
}
