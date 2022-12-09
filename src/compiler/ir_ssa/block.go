package ir

import (
	"fmt"
	"strings"
)

// Block 代码块
type Block struct {
	Depends map[*Block]struct{}
	Belong  *Function
	From    []*Block
	No      uint
	Stmts   []Stmt
}

func (self Block) String() string {
	var buf strings.Builder
	buf.WriteString(self.GetName())
	buf.WriteString(":\n")
	for _, s := range self.Stmts {
		buf.WriteString("  ")
		buf.WriteString(s.String())
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (self Block) GetName() string {
	return fmt.Sprintf("b%d", self.No)
}

// 获取本地变量当前计数
func (self *Block) getVarCount() uint {
	self.Belong.varCount++
	return self.Belong.varCount
}

// AddDepends 增加依赖
func (self *Block) AddDepends(b ...*Block) {
	for _, bb := range b {
		if bb == nil {
			continue
		}
		self.Depends[bb] = struct{}{}
	}
}

// NewReturn 新建函数返回
func (self *Block) NewReturn(v Value) *Return {
	if v != nil {
		self.AddDepends(v.GetBelong())
	}
	stmt := &Return{
		Belong: self,
		Value:  v,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewAdd 新建加
func (self *Block) NewAdd(l, r Value) *Add {
	_ = l.GetType().(NumberType)
	_ = r.GetType().(NumberType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Add{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewSub 新建减
func (self *Block) NewSub(l, r Value) *Sub {
	_ = l.GetType().(NumberType)
	_ = r.GetType().(NumberType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Sub{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewMul 新建乘
func (self *Block) NewMul(l, r Value) *Mul {
	_ = l.GetType().(NumberType)
	_ = r.GetType().(NumberType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Mul{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewDiv 新建除
func (self *Block) NewDiv(l, r Value) *Div {
	_ = l.GetType().(NumberType)
	_ = r.GetType().(NumberType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Div{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewMod 新建取余
func (self *Block) NewMod(l, r Value) *Mod {
	_ = l.GetType().(NumberType)
	_ = r.GetType().(NumberType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Mod{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewAnd 新建与
func (self *Block) NewAnd(l, r Value) *And {
	_ = l.GetType().(IntType)
	_ = r.GetType().(IntType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &And{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewOr 新建或
func (self *Block) NewOr(l, r Value) *Or {
	_ = l.GetType().(IntType)
	_ = r.GetType().(IntType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Or{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewXor 新建异或
func (self *Block) NewXor(l, r Value) *Xor {
	_ = l.GetType().(IntType)
	_ = r.GetType().(IntType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Xor{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewShl 新建左移
func (self *Block) NewShl(l, r Value) *Shl {
	_ = l.GetType().(IntType)
	_ = r.GetType().(IntType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Shl{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewShr 新建右移
func (self *Block) NewShr(l, r Value) *Shr {
	_ = l.GetType().(IntType)
	_ = r.GetType().(IntType)

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Shr{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewAlloc 新建栈分配
func (self *Block) NewAlloc(t Type) *Alloc {
	stmt := &Alloc{
		Belong: self,
		No:     self.getVarCount(),
		Type:   t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewStore 新建赋值
func (self *Block) NewStore(f, t Value) *Store {
	ptr := t.GetType().(*TypePtr)
	if !f.GetType().Equal(ptr.Elem) {
		panic("")
	}

	self.AddDepends(f.GetBelong(), t.GetBelong())
	stmt := &Store{
		Belong: self,
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewLoad 新建取值
func (self *Block) NewLoad(p Value) *Load {
	_ = p.GetType().(*TypePtr)

	self.AddDepends(p.GetBelong())
	stmt := &Load{
		Belong: self,
		No:     self.getVarCount(),
		Ptr:    p,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewCall 新建调用
func (self *Block) NewCall(f Value, args ...Value) *Call {
	ret := f.GetType().(*TypePtr).Elem.(*TypeFunc).Ret

	self.AddDepends(f.GetBelong())
	for _, a := range args {
		self.AddDepends(a.GetBelong())
	}
	var no uint
	if !ret.Equal(None) {
		no = self.getVarCount()
	}
	stmt := &Call{
		Belong: self,
		No:     no,
		Func:   f,
		Args:   args,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewArrayIndex 新建数组索引
func (self *Block) NewArrayIndex(f, i Value) *ArrayIndex {
	ft := f.GetType()
	if (!IsTypePtr(ft) || !IsArrayType(ft.(*TypePtr).Elem)) && !IsArrayType(ft) {
		panic("")
	}
	if !i.GetType().Equal(Usize) {
		panic("")
	}

	self.AddDepends(f.GetBelong(), i.GetBelong())
	stmt := &ArrayIndex{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		Index:  i,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewPtrIndex 新建指针索引
func (self *Block) NewPtrIndex(f, i Value) *PtrIndex {
	ft := f.GetType()
	if !IsTypePtr(ft) {
		panic("")
	}
	if !i.GetType().Equal(Usize) {
		panic("")
	}

	self.AddDepends(f.GetBelong(), i.GetBelong())
	stmt := &PtrIndex{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		Index:  i,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewEq 新建相等
func (self *Block) NewEq(l, r Value) *Eq {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) && !IsTypePtr(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Eq{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewNe 新建不等
func (self *Block) NewNe(l, r Value) *Ne {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) && !IsTypePtr(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Ne{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewLt 新建小于
func (self *Block) NewLt(l, r Value) *Lt {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Lt{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewLe 新建小于等于
func (self *Block) NewLe(l, r Value) *Le {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Le{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewGt 新建大于
func (self *Block) NewGt(l, r Value) *Gt {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Gt{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewGe 新建大于等于
func (self *Block) NewGe(l, r Value) *Ge {
	lt, rt := l.GetType(), r.GetType()
	if !lt.Equal(rt) {
		panic("")
	} else if !IsNumberType(lt) {
		panic("")
	}

	self.AddDepends(l.GetBelong(), r.GetBelong())
	stmt := &Ge{
		Belong: self,
		No:     self.getVarCount(),
		Left:   l,
		Right:  r,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewGoto 新建跳转
func (self *Block) NewGoto(b *Block) *Goto {
	stmt := &Goto{
		Belong: self,
		Target: b,
	}
	self.Stmts = append(self.Stmts, stmt)
	b.From = append(b.From, self)
	return stmt
}

// NewCondGoto 新建条件跳转
func (self *Block) NewCondGoto(c Value, t, f *Block) *CondGoto {
	if !c.GetType().Equal(I8) {
		panic("")
	}
	self.AddDepends(c.GetBelong())
	stmt := &CondGoto{
		Belong: self,
		Cond:   c,
		True:   t,
		False:  f,
	}
	self.Stmts = append(self.Stmts, stmt)
	t.From = append(t.From, self)
	f.From = append(f.From, self)
	return stmt
}

// NewUnreachable 新建不可达
func (self *Block) NewUnreachable() *Unreachable {
	stmt := &Unreachable{
		Belong: self,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewPhi 新建跳转获取
func (self *Block) NewPhi(blocks []*Block, values []Value) *Phi {
	if len(blocks) != len(values) || len(blocks) == 0 {
		panic("")
	}
	for _, v := range values {
		self.AddDepends(v.GetBelong())
	}
	stmt := &Phi{
		Belong: self,
		No:     self.getVarCount(),
		Froms:  blocks,
		Values: values,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewStructIndex 新建结构体索引
func (self *Block) NewStructIndex(f Value, i uint) *StructIndex {
	ft := f.GetType()
	if (!IsTypePtr(ft) || !IsStructType(ft.(*TypePtr).Elem)) && !IsStructType(ft) {
		panic("")
	}

	self.AddDepends(f.GetBelong())
	stmt := &StructIndex{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		Index:  i,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewItoi 新建整型转换
func (self *Block) NewItoi(f Value, t IntType) *Itoi {
	_ = f.GetType().(IntType)

	self.AddDepends(f.GetBelong())
	stmt := &Itoi{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewFtof 新建浮点型转换
func (self *Block) NewFtof(f Value, t *TypeFloat) *Ftof {
	_ = f.GetType().(*TypeFloat)

	self.AddDepends(f.GetBelong())
	stmt := &Ftof{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewItof 新建整型转浮点型
func (self *Block) NewItof(f Value, t *TypeFloat) *Itof {
	_ = f.GetType().(IntType)

	self.AddDepends(f.GetBelong())
	stmt := &Itof{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewFtoi 新建浮点型转整型
func (self *Block) NewFtoi(f Value, t IntType) *Ftoi {
	_ = f.GetType().(*TypeFloat)

	self.AddDepends(f.GetBelong())
	stmt := &Ftoi{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewPtop 新建指针转指针
func (self *Block) NewPtop(f Value, t *TypePtr) *Ptop {
	_ = f.GetType().(*TypePtr)

	self.AddDepends(f.GetBelong())
	stmt := &Ptop{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewPtoi 新建指针转整型
func (self *Block) NewPtoi(f Value, t IntType) *Ptoi {
	_ = f.GetType().(*TypePtr)

	self.AddDepends(f.GetBelong())
	stmt := &Ptoi{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}

// NewItop 新建整型转指针
func (self *Block) NewItop(f Value, t *TypePtr) *Itop {
	_ = f.GetType().(IntType)

	self.AddDepends(f.GetBelong())
	stmt := &Itop{
		Belong: self,
		No:     self.getVarCount(),
		From:   f,
		To:     t,
	}
	self.Stmts = append(self.Stmts, stmt)
	return stmt
}
