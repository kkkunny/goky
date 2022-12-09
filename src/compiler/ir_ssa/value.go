package ir

import (
	"fmt"
	"strconv"
	"strings"
)

// Value 值
type Value interface {
	GetName() string
	GetType() Type
	GetBelong() *Block
}

// Constant 常量
type Constant interface {
	Value
	IsZero() bool
}

// Int 整数
type Int struct {
	Type  IntType
	Value int64
}

// NewInt 新建整数
func NewInt(t IntType, v int64) *Int {
	return &Int{
		Type:  t,
		Value: v,
	}
}

func (self Int) GetName() string {
	return strconv.FormatInt(self.Value, 10)
}

func (self Int) GetType() Type {
	return self.Type
}

func (self Int) IsZero() bool {
	return self.Value == 0
}

func (self Int) GetBelong() *Block {
	return nil
}

// Float 浮点数
type Float struct {
	Type  *TypeFloat
	Value float64
}

// NewFloat 新建浮点数
func NewFloat(t *TypeFloat, v float64) *Float {
	return &Float{
		Type:  t,
		Value: v,
	}
}

func (self Float) GetName() string {
	return strconv.FormatFloat(self.Value, 'f', -1, 64)
}

func (self Float) GetType() Type {
	return self.Type
}

func (self Float) IsZero() bool {
	return self.Value == 0
}

func (self Float) GetBelong() *Block {
	return nil
}

// Empty 空
type Empty struct {
	Type Type
}

// NewEmpty 新建空
func NewEmpty(t Type) *Empty {
	if !IsArrayType(t) && !IsFuncType(t) && !IsTypePtr(t) && !IsStructType(t) {
		panic("")
	}
	return &Empty{
		Type: t,
	}
}

func (self Empty) GetName() string {
	return "empty"
}

func (self Empty) GetType() Type {
	return self.Type
}

func (self Empty) IsZero() bool {
	return true
}

func (self Empty) GetBelong() *Block {
	return nil
}

// Array 数组
type Array struct {
	Elems []Constant
}

// NewArray 新建数组
func NewArray(elem ...Constant) *Array {
	if len(elem) == 0 {
		panic("")
	}
	return &Array{Elems: elem}
}

func (self Array) GetName() string {
	elems := make([]string, len(self.Elems))
	for i, e := range self.Elems {
		elems[i] = e.GetName()
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

func (self Array) GetType() Type {
	return NewArrayType(uint(len(self.Elems)), self.Elems[0].GetType())
}

func (self Array) IsZero() bool {
	for _, e := range self.Elems {
		if !e.IsZero() {
			return false
		}
	}
	return true
}

func (self Array) GetBelong() *Block {
	return nil
}

// Struct 结构体
type Struct struct {
	Elems []Constant
}

// NewStruct 新建结构体
func NewStruct(elem ...Constant) *Struct {
	if len(elem) == 0 {
		panic("")
	}
	return &Struct{Elems: elem}
}

func (self Struct) GetName() string {
	var elems []string
	for i, e := range self.Elems {
		elems[i] = e.GetName()
	}
	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}

func (self Struct) GetType() Type {
	elems := make([]Type, len(self.Elems))
	for i, e := range self.Elems {
		elems[i] = e.GetType()
	}
	return NewStructType(elems...)
}

func (self Struct) IsZero() bool {
	for _, e := range self.Elems {
		if !e.IsZero() {
			return false
		}
	}
	return true
}

func (self Struct) GetBelong() *Block {
	return nil
}
