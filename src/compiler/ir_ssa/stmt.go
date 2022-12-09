package ir

import (
	"fmt"
	"strings"
)

type Stmt interface {
	fmt.Stringer
	stmt()
}

type Return struct {
	Belong *Block
	Value  Value
}

func (self Return) stmt() {}

func (self Return) String() string {
	if self.Value == nil {
		return "ret"
	} else {
		return fmt.Sprintf("ret %s", self.Value.GetName())
	}
}

type Add struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Add) stmt() {}

func (self Add) String() string {
	return fmt.Sprintf("%s %s = add %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Add) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Add) GetType() Type {
	return self.Left.GetType()
}

func (self Add) GetBelong() *Block {
	return self.Belong
}

type Sub struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Sub) stmt() {}

func (self Sub) String() string {
	return fmt.Sprintf("%s %s = sub %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Sub) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Sub) GetType() Type {
	return self.Left.GetType()
}

func (self Sub) GetBelong() *Block {
	return self.Belong
}

type Mul struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Mul) stmt() {}

func (self Mul) String() string {
	return fmt.Sprintf("%s %s = mul %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Mul) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Mul) GetType() Type {
	return self.Left.GetType()
}

func (self Mul) GetBelong() *Block {
	return self.Belong
}

type Div struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Div) stmt() {}

func (self Div) String() string {
	return fmt.Sprintf("%s %s = div %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Div) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Div) GetType() Type {
	return self.Left.GetType()
}

func (self Div) GetBelong() *Block {
	return self.Belong
}

type Mod struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Mod) stmt() {}

func (self Mod) String() string {
	return fmt.Sprintf("%s %s = mod %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Mod) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Mod) GetType() Type {
	return self.Left.GetType()
}

func (self Mod) GetBelong() *Block {
	return self.Belong
}

type And struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self And) stmt() {}

func (self And) String() string {
	return fmt.Sprintf("%s %s = and %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self And) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self And) GetType() Type {
	return self.Left.GetType()
}

func (self And) GetBelong() *Block {
	return self.Belong
}

type Or struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Or) stmt() {}

func (self Or) String() string {
	return fmt.Sprintf("%s %s = or %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Or) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Or) GetType() Type {
	return self.Left.GetType()
}

func (self Or) GetBelong() *Block {
	return self.Belong
}

type Xor struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Xor) stmt() {}

func (self Xor) String() string {
	return fmt.Sprintf("%s %s = xor %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Xor) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Xor) GetType() Type {
	return self.Left.GetType()
}

func (self Xor) GetBelong() *Block {
	return self.Belong
}

type Shl struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Shl) stmt() {}

func (self Shl) String() string {
	return fmt.Sprintf("%s %s = shl %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Shl) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Shl) GetType() Type {
	return self.Left.GetType()
}

func (self Shl) GetBelong() *Block {
	return self.Belong
}

type Shr struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Shr) stmt() {}

func (self Shr) String() string {
	return fmt.Sprintf("%s %s = shr %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Shr) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Shr) GetType() Type {
	return self.Left.GetType()
}

func (self Shr) GetBelong() *Block {
	return self.Belong
}

// Alloc 栈分配
type Alloc struct {
	Belong *Block
	No     uint
	Type   Type
}

func (self Alloc) stmt() {}

func (self Alloc) String() string {
	return fmt.Sprintf("%s %s = alloc %s", self.GetType(), self.GetName(), self.Type)
}

func (self Alloc) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Alloc) GetType() Type {
	return NewPtrType(self.Type)
}

func (self Alloc) GetBelong() *Block {
	return self.Belong
}

// Store 赋值
type Store struct {
	Belong   *Block
	From, To Value
}

func (self Store) stmt() {}

func (self Store) String() string {
	return fmt.Sprintf("store %s to %s", self.From.GetName(), self.To.GetName())
}

// Load 取值
type Load struct {
	Belong *Block
	No     uint
	Ptr    Value
}

func (self Load) stmt() {}

func (self Load) String() string {
	return fmt.Sprintf("%s %s = load %s", self.GetType(), self.GetName(), self.Ptr.GetName())
}

func (self Load) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Load) GetType() Type {
	return self.Ptr.GetType().(*TypePtr).Elem
}

func (self Load) GetBelong() *Block {
	return self.Belong
}

// Call 调用
type Call struct {
	Belong *Block
	No     uint
	Func   Value
	Args   []Value
}

func (self Call) stmt() {}

func (self Call) String() string {
	args := make([]string, len(self.Args))
	for i, a := range self.Args {
		args[i] = a.GetName()
	}
	if self.GetType().Equal(None) {
		return fmt.Sprintf("call %s(%s)", self.Func.GetName(), strings.Join(args, ", "))
	} else {
		return fmt.Sprintf("%s %s = call %s(%s)", self.GetType(), self.GetName(), self.Func.GetName(), strings.Join(args, ", "))
	}
}

func (self Call) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Call) GetType() Type {
	return self.Func.GetType().(*TypePtr).Elem.(*TypeFunc).Ret
}

func (self Call) GetBelong() *Block {
	return self.Belong
}

// ArrayIndex 数组索引
type ArrayIndex struct {
	Belong *Block
	No     uint
	From   Value
	Index  Value
}

func (self ArrayIndex) stmt() {}

func (self ArrayIndex) String() string {
	return fmt.Sprintf("%s %s = index %s, %s", self.GetType(), self.GetName(), self.From.GetName(), self.Index.GetName())
}

func (self ArrayIndex) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self ArrayIndex) GetType() Type {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return NewPtrType(typ.Elem.(*TypeArray).Elem)
	case *TypeArray:
		return typ.Elem
	default:
		panic("")
	}
}

func (self ArrayIndex) GetBelong() *Block {
	return self.Belong
}

func (self ArrayIndex) GetElemType() Type {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return typ.Elem.(*TypeArray).Elem
	case *TypeArray:
		return typ.Elem
	default:
		panic("")
	}
}

func (self ArrayIndex) GetFromType() *TypeArray {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return typ.Elem.(*TypeArray)
	case *TypeArray:
		return typ
	default:
		panic("")
	}
}

// PtrIndex 指针索引
type PtrIndex struct {
	Belong *Block
	No     uint
	From   Value
	Index  Value
}

func (self PtrIndex) stmt() {}

func (self PtrIndex) String() string {
	return fmt.Sprintf("%s %s = index %s, %s", self.GetType(), self.GetName(), self.From.GetName(), self.Index.GetName())
}

func (self PtrIndex) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self PtrIndex) GetType() Type {
	return NewPtrType(self.From.GetType().(*TypePtr).Elem)
}

func (self PtrIndex) GetBelong() *Block {
	return self.Belong
}

func (self PtrIndex) GetElemType() Type {
	return self.From.GetType().(*TypePtr).Elem
}

func (self PtrIndex) GetFromType() *TypePtr {
	return self.From.GetType().(*TypePtr)
}

type Eq struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Eq) stmt() {}

func (self Eq) String() string {
	return fmt.Sprintf("%s %s = eq %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Eq) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Eq) GetType() Type {
	return I8
}

func (self Eq) GetBelong() *Block {
	return self.Belong
}

type Ne struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Ne) stmt() {}

func (self Ne) String() string {
	return fmt.Sprintf("%s %s = ne %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Ne) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ne) GetType() Type {
	return I8
}

func (self Ne) GetBelong() *Block {
	return self.Belong
}

type Lt struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Lt) stmt() {}

func (self Lt) String() string {
	return fmt.Sprintf("%s %s = lt %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Lt) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Lt) GetType() Type {
	return I8
}

func (self Lt) GetBelong() *Block {
	return self.Belong
}

type Le struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Le) stmt() {}

func (self Le) String() string {
	return fmt.Sprintf("%s %s = le %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Le) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Le) GetType() Type {
	return I8
}

func (self Le) GetBelong() *Block {
	return self.Belong
}

type Gt struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Gt) stmt() {}

func (self Gt) String() string {
	return fmt.Sprintf("%s %s = gt %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Gt) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Gt) GetType() Type {
	return I8
}

func (self Gt) GetBelong() *Block {
	return self.Belong
}

type Ge struct {
	Belong      *Block
	No          uint
	Left, Right Value
}

func (self Ge) stmt() {}

func (self Ge) String() string {
	return fmt.Sprintf("%s %s = ge %s, %s", self.GetType(), self.GetName(), self.Left.GetName(), self.Right.GetName())
}

func (self Ge) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ge) GetType() Type {
	return I8
}

func (self Ge) GetBelong() *Block {
	return self.Belong
}

type Goto struct {
	Belong *Block
	Target *Block
}

func (self Goto) stmt() {}

func (self Goto) String() string {
	return fmt.Sprintf("goto %s", self.Target.GetName())
}

type CondGoto struct {
	Belong      *Block
	Cond        Value
	True, False *Block
}

func (self CondGoto) stmt() {}

func (self CondGoto) String() string {
	return fmt.Sprintf("if %s goto %s or %s", self.Cond.GetName(), self.True.GetName(), self.False.GetName())
}

type Unreachable struct {
	Belong *Block
}

func (self Unreachable) stmt() {}

func (self Unreachable) String() string {
	return fmt.Sprintf("unreachable")
}

type Phi struct {
	Belong *Block
	No     uint
	Froms  []*Block
	Values []Value
}

func (self Phi) stmt() {}

func (self Phi) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s %s = phi ", self.GetType(), self.GetName()))
	for i, b := range self.Froms {
		buf.WriteString(fmt.Sprintf("%s:%s", b.GetName(), self.Values[i].GetName()))
		if i < len(self.Froms)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func (self Phi) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Phi) GetType() Type {
	return self.Values[0].GetType()
}

func (self Phi) GetBelong() *Block {
	return self.Belong
}

// StructIndex 结构体索引
type StructIndex struct {
	Belong *Block
	No     uint
	From   Value
	Index  uint
}

func (self StructIndex) stmt() {}

func (self StructIndex) String() string {
	return fmt.Sprintf("%s %s = index %s, %d", self.GetType(), self.GetName(), self.From.GetName(), self.Index)
}

func (self StructIndex) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self StructIndex) GetType() Type {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return NewPtrType(typ.Elem.(*TypeStruct).Elems[self.Index])
	case *TypeStruct:
		return typ.Elems[self.Index]
	default:
		panic("")
	}
}

func (self StructIndex) GetBelong() *Block {
	return self.Belong
}

func (self StructIndex) GetFromType() *TypeStruct {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return typ.Elem.(*TypeStruct)
	case *TypeStruct:
		return typ
	default:
		panic("")
	}
}

func (self StructIndex) GetElemType() Type {
	switch typ := self.From.GetType().(type) {
	case *TypePtr:
		return typ.Elem.(*TypeStruct).Elems[self.Index]
	case *TypeStruct:
		return typ.Elems[self.Index]
	default:
		panic("")
	}
}

// Itoi 整型转换
type Itoi struct {
	Belong *Block
	No     uint
	From   Value
	To     IntType
}

func (self Itoi) stmt() {}

func (self Itoi) String() string {
	return fmt.Sprintf("%s %s = itoi %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Itoi) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Itoi) GetType() Type {
	return self.To
}

func (self Itoi) GetBelong() *Block {
	return self.Belong
}

// Ftof 浮点型转换
type Ftof struct {
	Belong *Block
	No     uint
	From   Value
	To     *TypeFloat
}

func (self Ftof) stmt() {}

func (self Ftof) String() string {
	return fmt.Sprintf("%s %s = ftof %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Ftof) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ftof) GetType() Type {
	return self.To
}

func (self Ftof) GetBelong() *Block {
	return self.Belong
}

// Itof 整型转浮点型
type Itof struct {
	Belong *Block
	No     uint
	From   Value
	To     *TypeFloat
}

func (self Itof) stmt() {}

func (self Itof) String() string {
	return fmt.Sprintf("%s %s = itof %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Itof) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Itof) GetType() Type {
	return self.To
}

func (self Itof) GetBelong() *Block {
	return self.Belong
}

// Ftoi 浮点型转整型
type Ftoi struct {
	Belong *Block
	No     uint
	From   Value
	To     IntType
}

func (self Ftoi) stmt() {}

func (self Ftoi) String() string {
	return fmt.Sprintf("%s %s = ftoi %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Ftoi) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ftoi) GetType() Type {
	return self.To
}

func (self Ftoi) GetBelong() *Block {
	return self.Belong
}

// Ptop 指针转指针
type Ptop struct {
	Belong *Block
	No     uint
	From   Value
	To     *TypePtr
}

func (self Ptop) stmt() {}

func (self Ptop) String() string {
	return fmt.Sprintf("%s %s = ptop %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Ptop) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ptop) GetType() Type {
	return self.To
}

func (self Ptop) GetBelong() *Block {
	return self.Belong
}

// Ptoi 指针转整型
type Ptoi struct {
	Belong *Block
	No     uint
	From   Value
	To     IntType
}

func (self Ptoi) stmt() {}

func (self Ptoi) String() string {
	return fmt.Sprintf("%s %s = ptoi %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Ptoi) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Ptoi) GetType() Type {
	return self.To
}

func (self Ptoi) GetBelong() *Block {
	return self.Belong
}

// Itop 整型转指针
type Itop struct {
	Belong *Block
	No     uint
	From   Value
	To     *TypePtr
}

func (self Itop) stmt() {}

func (self Itop) String() string {
	return fmt.Sprintf("%s %s = itop %s", self.GetType(), self.GetName(), self.From.GetName())
}

func (self Itop) GetName() string {
	return fmt.Sprintf("v%d", self.No)
}

func (self Itop) GetType() Type {
	return self.To
}

func (self Itop) GetBelong() *Block {
	return self.Belong
}
