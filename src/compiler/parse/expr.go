package parse

import (
	"github.com/kkkunny/klang/src/compiler/lex"
	"github.com/kkkunny/klang/src/compiler/utils"
	"strconv"
)

// Expr 表达式
type Expr interface {
	Stmt
	Expr()
}

// Call 调用
type Call struct {
	Pos  utils.Position
	Func Expr
	Args []Expr
}

func NewCall(pos utils.Position, f Expr, a ...Expr) *Call {
	return &Call{
		Pos:  pos,
		Func: f,
		Args: a,
	}
}

func (self Call) Position() utils.Position {
	return self.Pos
}

func (self Call) Stmt() {}

func (self Call) Expr() {}

// Int 整数
type Int struct {
	Token lex.Token
	Value int64
}

func NewInt(tok lex.Token, v int64) *Int {
	return &Int{
		Token: tok,
		Value: v,
	}
}

func (self Int) Position() utils.Position {
	return self.Token.Pos
}

func (self Int) Stmt() {}

func (self Int) Expr() {}

// Float 浮点数
type Float struct {
	Token lex.Token
	Value float64
}

func NewFloat(tok lex.Token, v float64) *Float {
	return &Float{
		Token: tok,
		Value: v,
	}
}

func (self Float) Position() utils.Position {
	return self.Token.Pos
}

func (self Float) Stmt() {}

func (self Float) Expr() {}

// Bool 布尔
type Bool struct {
	Token lex.Token
	Value bool
}

func NewBool(tok lex.Token, v bool) *Bool {
	return &Bool{
		Token: tok,
		Value: v,
	}
}

func (self Bool) Position() utils.Position {
	return self.Token.Pos
}

func (self Bool) Stmt() {}

func (self Bool) Expr() {}

// Char 字符
type Char struct {
	Token lex.Token
	Value rune
}

func NewChar(tok lex.Token, v rune) *Char {
	return &Char{
		Token: tok,
		Value: v,
	}
}

func (self Char) Position() utils.Position {
	return self.Token.Pos
}

func (self Char) Stmt() {}

func (self Char) Expr() {}

// String 字符串
type String struct {
	Token lex.Token
	Value []rune
}

func NewString(tok lex.Token, v []rune) *String {
	return &String{
		Token: tok,
		Value: v,
	}
}

func (self String) Position() utils.Position {
	return self.Token.Pos
}

func (self String) Stmt() {}

func (self String) Expr() {}

// CString 字符串
type CString struct {
	Token lex.Token
	Value []byte
}

func NewCString(tok lex.Token, v []byte) *CString {
	return &CString{
		Token: tok,
		Value: v,
	}
}

func (self CString) Position() utils.Position {
	return self.Token.Pos
}

func (self CString) Stmt() {}

func (self CString) Expr() {}

// Null 空指针
type Null struct {
	Token lex.Token
}

func NewNull(tok lex.Token) *Null {
	return &Null{Token: tok}
}

func (self Null) Position() utils.Position {
	return self.Token.Pos
}

func (self Null) Stmt() {}

func (self Null) Expr() {}

// Ident 标识符
type Ident struct {
	Pkg  *lex.Token
	Name lex.Token
}

func NewIdent(pkg *lex.Token, name lex.Token) *Ident {
	return &Ident{
		Pkg:  pkg,
		Name: name,
	}
}

func (self Ident) Position() utils.Position {
	if self.Pkg != nil {
		return utils.MixPosition(self.Pkg.Pos, self.Name.Pos)
	} else {
		return self.Name.Pos
	}
}

func (self Ident) Stmt() {}

func (self Ident) Expr() {}

// Array 数组
type Array struct {
	Pos   utils.Position
	Elems []Expr
}

func NewArray(pos utils.Position, elem ...Expr) *Array {
	return &Array{
		Pos:   pos,
		Elems: elem,
	}
}

func (self Array) Position() utils.Position {
	return self.Pos
}

func (self Array) Stmt() {}

func (self Array) Expr() {}

// TupleOrExpr 元组或者括号表达式
type TupleOrExpr struct {
	Pos   utils.Position
	Elems []Expr
}

func NewTuple(pos utils.Position, elem ...Expr) *TupleOrExpr {
	return &TupleOrExpr{
		Pos:   pos,
		Elems: elem,
	}
}

func (self TupleOrExpr) Position() utils.Position {
	return self.Pos
}

func (self TupleOrExpr) Stmt() {}

func (self TupleOrExpr) Expr() {}

// Struct 结构体
type Struct struct {
	Pos    utils.Position
	Fields []Expr
}

func NewStruct(pos utils.Position, field ...Expr) *Struct {
	return &Struct{
		Pos:    pos,
		Fields: field,
	}
}

func (self Struct) Position() utils.Position {
	return self.Pos
}

func (self Struct) Stmt() {}

func (self Struct) Expr() {}

// Unary 一元表达式
type Unary struct {
	Opera lex.Token
	Value Expr
}

func NewUnary(op lex.Token, v Expr) *Unary {
	return &Unary{
		Opera: op,
		Value: v,
	}
}

func (self Unary) Position() utils.Position {
	return utils.MixPosition(self.Opera.Pos, self.Value.Position())
}

func (self Unary) Stmt() {}

func (self Unary) Expr() {}

// Dot 点
type Dot struct {
	Front Expr
	End   lex.Token
}

func NewDot(f Expr, e lex.Token) *Dot {
	return &Dot{
		Front: f,
		End:   e,
	}
}

func (self Dot) Position() utils.Position {
	return utils.MixPosition(self.Front.Position(), self.End.Pos)
}

func (self Dot) Stmt() {}

func (self Dot) Expr() {}

// Index 索引
type Index struct {
	Pos   utils.Position
	Front Expr
	Index Expr
}

func NewIndex(pos utils.Position, f, e Expr) *Index {
	return &Index{
		Pos:   pos,
		Front: f,
		Index: e,
	}
}

func (self Index) Position() utils.Position {
	return self.Pos
}

func (self Index) Stmt() {}

func (self Index) Expr() {}

// Covert 类型转换
type Covert struct {
	From Expr
	To   Type
}

func NewCovert(f Expr, t Type) *Covert {
	return &Covert{
		From: f,
		To:   t,
	}
}

func (self Covert) Position() utils.Position {
	return utils.MixPosition(self.From.Position(), self.To.Position())
}

func (self Covert) Stmt() {}

func (self Covert) Expr() {}

// Ternary 三元表达式
type Ternary struct {
	Cond, True, False Expr
}

func NewTernary(cond, t, f Expr) *Ternary {
	return &Ternary{
		Cond:  cond,
		True:  t,
		False: f,
	}
}

func (self Ternary) Position() utils.Position {
	return utils.MixPosition(self.Cond.Position(), self.False.Position())
}

func (self Ternary) Stmt() {}

func (self Ternary) Expr() {}

// Binary 二元表达式
type Binary struct {
	Opera       lex.Token
	Left, Right Expr
}

func NewBinary(op lex.Token, l, r Expr) *Binary {
	return &Binary{
		Opera: op,
		Left:  l,
		Right: r,
	}
}

func (self Binary) Position() utils.Position {
	return utils.MixPosition(self.Left.Position(), self.Right.Position())
}

func (self Binary) Stmt() {}

func (self Binary) Expr() {}

// ****************************************************************

// 表达式
func (self *Parser) parseExpr() Expr {
	return self.parseBinaryExpr(0)
}

// 表达式列表（至少一个）
func (self *Parser) parseExprListAtLeastOne(sep lex.TokenKind) (toks []Expr) {
	for {
		toks = append(toks, self.parseExpr())
		if !self.skipNextIs(sep) {
			break
		}
	}
	return toks
}

// 单表达式
func (self *Parser) parsePrimaryExpr() Expr {
	switch self.nextTok.Kind {
	case lex.INT:
		self.next()
		v, err := strconv.ParseInt(self.curTok.Source, 10, 64)
		if err != nil {
			self.throwErrorf(self.curTok.Pos, "out of integer size")
		}
		return NewInt(self.curTok, v)
	case lex.FLOAT:
		self.next()
		v, err := strconv.ParseFloat(self.curTok.Source, 64)
		if err != nil {
			self.throwErrorf(self.curTok.Pos, "out of float size")
		}
		return NewFloat(self.curTok, v)
	case lex.TRUE:
		self.next()
		return NewBool(self.curTok, true)
	case lex.FALSE:
		self.next()
		return NewBool(self.curTok, false)
	case lex.CHAR:
		self.next()
		return NewChar(self.curTok, []rune(self.curTok.Source)[1])
	case lex.STRING:
		self.next()
		runes := []rune(self.curTok.Source[1 : len(self.curTok.Source)-1])
		return NewString(self.curTok, runes)
	case lex.CSTRING:
		self.next()
		bytes := []byte(self.curTok.Source[1 : len(self.curTok.Source)-1])
		return NewCString(self.curTok, bytes)
	case lex.NULL:
		self.next()
		return NewNull(self.curTok)
	case lex.IDENT:
		self.next()
		name := self.curTok
		var pkg *lex.Token
		if self.skipNextIs(lex.CLL) {
			pkg = &name
			name = self.expectNextIs(lex.IDENT)
		}
		return NewIdent(pkg, name)
	case lex.LPA:
		self.next()
		begin := self.curTok.Pos
		var elems []Expr
		if !self.nextIs(lex.RPA) {
			elems = self.parseExprListAtLeastOne(lex.COM)
		}
		end := self.expectNextIs(lex.RPA).Pos
		return NewArray(utils.MixPosition(begin, end), elems...)
	case lex.LBA:
		self.next()
		begin := self.curTok.Pos
		var elems []Expr
		if !self.nextIs(lex.RBA) {
			elems = self.parseExprListAtLeastOne(lex.COM)
		}
		end := self.expectNextIs(lex.RBA).Pos
		return NewTuple(utils.MixPosition(begin, end), elems...)
	case lex.LBR:
		self.next()
		begin := self.curTok.Pos
		var fields []Expr
		if !self.nextIs(lex.RBR) {
			fields = self.parseExprListAtLeastOne(lex.COM)
		}
		end := self.expectNextIs(lex.RBR).Pos
		return NewStruct(utils.MixPosition(begin, end), fields...)
	default:
		self.throwErrorf(self.nextTok.Pos, "unknown expression")
		return nil
	}
}

// 一元表达式前缀
func (self *Parser) parsePrefixUnaryExpr() Expr {
	switch self.nextTok.Kind {
	case lex.SUB, lex.NEG, lex.NOT, lex.AND, lex.MUL:
		self.next()
		op := self.curTok
		v := self.parsePrefixUnaryExpr()
		return NewUnary(op, v)
	default:
		return self.parsePrimaryExpr()
	}
}

// 一元表达式后缀
func (self *Parser) parseSuffixUnaryExpr(front Expr) Expr {
	switch self.nextTok.Kind {
	case lex.DOT:
		self.next()
		end := self.expectNextIs(lex.IDENT)
		front = NewDot(front, end)
	case lex.LPA:
		self.next()
		var args []Expr
		if !self.nextIs(lex.RPA) {
			args = self.parseExprListAtLeastOne(lex.COM)
		}
		end := self.expectNextIs(lex.RPA).Pos
		front = NewCall(utils.MixPosition(front.Position(), end), front, args...)
	case lex.LBA:
		self.next()
		self.expectNextIs(lex.LBA)
		index := self.parseExpr()
		end := self.expectNextIs(lex.RBA).Pos
		front = NewIndex(utils.MixPosition(front.Position(), end), front, index)
	default:
		return front
	}
	return self.parseSuffixUnaryExpr(front)
}

// 一元表达式末尾
func (self *Parser) parseUnaryTailExpr(front Expr) Expr {
	switch self.nextTok.Kind {
	case lex.AS:
		self.next()
		t := self.parseType()
		front = NewCovert(front, t)
	case lex.QUO:
		self.next()
		t := self.parseUnaryExpr()
		self.expectNextIs(lex.COL)
		f := self.parseUnaryExpr()
		front = NewTernary(front, t, f)
	default:
		return front
	}
	return self.parseUnaryTailExpr(front)
}

// 一元表达式
func (self *Parser) parseUnaryExpr() Expr {
	front := self.parsePrefixUnaryExpr()
	front = self.parseSuffixUnaryExpr(front)
	return self.parseUnaryTailExpr(front)
}

// 二元表达式
func (self *Parser) parseBinaryExpr(prior uint8) Expr {
	left := self.parseUnaryExpr()

	for {
		nextOpera := self.nextTok
		if prior >= nextOpera.Kind.Priority() {
			break
		}
		self.next()
		right := self.parseBinaryExpr(nextOpera.Kind.Priority())
		left = &Binary{
			Opera: nextOpera,
			Left:  left,
			Right: right,
		}
	}

	return left
}
