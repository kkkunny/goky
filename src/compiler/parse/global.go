package parse

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/lex"
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Global 全局
type Global interface {
	Ast
	Global()
}

// Import 包导入
type Import struct {
	Pos      utils.Position
	Packages []lex.Token
	Alias    *lex.Token
}

func NewImport(pos utils.Position, pkgs []lex.Token, alias *lex.Token) *Import {
	return &Import{
		Pos:      pos,
		Packages: pkgs,
		Alias:    alias,
	}
}

func (self Import) Position() utils.Position {
	return self.Pos
}

func (self Import) Global() {}

// TypeDef 类型定义
type TypeDef struct {
	Pos    utils.Position
	Public bool
	Name   lex.Token
	Target Type
}

func NewTypeDef(pos utils.Position, pub bool, name lex.Token, target Type) *TypeDef {
	return &TypeDef{
		Pos:    pos,
		Public: pub,
		Name:   name,
		Target: target,
	}
}

func (self TypeDef) Position() utils.Position {
	return self.Pos
}

func (self TypeDef) Global() {}

// Function 函数
type Function struct {
	Pos    utils.Position
	Attrs  []Attr
	Public bool
	Ret    Type
	Name   lex.Token
	Params []*NameOrNilAndType
	Body   *Block // 可能为空
}

func NewFunction(pos utils.Position, attrs []Attr, pub bool, ret Type, name lex.Token, params []*NameOrNilAndType, body *Block) *Function {
	return &Function{
		Pos:    pos,
		Attrs:  attrs,
		Public: pub,
		Ret:    ret,
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (self Function) Position() utils.Position {
	return self.Pos
}

func (self Function) Global() {}

// Method 方法
type Method struct {
	Pos    utils.Position
	Attrs  []Attr
	Public bool
	Self   lex.Token
	Ret    Type
	Name   lex.Token
	Params []*NameOrNilAndType
	Body   *Block
}

func NewMethod(pos utils.Position, attrs []Attr, pub bool, self lex.Token, ret Type, name lex.Token, params []*NameOrNilAndType, body *Block) *Method {
	return &Method{
		Pos:    pos,
		Attrs:  attrs,
		Public: pub,
		Self:   self,
		Ret:    ret,
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (self Method) Position() utils.Position {
	return self.Pos
}

func (self Method) Global() {}

// GlobalValue 全局变量
type GlobalValue struct {
	Attrs    []Attr
	Public   bool
	Variable *Variable
}

func NewGlobalValue(pos utils.Position, attrs []Attr, pub bool, t Type, name lex.Token, v Expr) *GlobalValue {
	return &GlobalValue{
		Attrs:    attrs,
		Public:   pub,
		Variable: NewVariable(pos, t, name, v),
	}
}

func (self GlobalValue) Position() utils.Position {
	return self.Variable.Pos
}

func (self GlobalValue) Global() {}

// ****************************************************************

var (
	errStrUnknownGlobal = "unknown global"
	errStrCanNotUseAttr = "can not use this attribute"
)

// 全局
func (self *Parser) parseGlobal() Global {
	var pub *lex.Token
	if self.skipNextIs(lex.PUB) {
		pub = &self.curTok
	}

	switch self.nextTok.Kind {
	case lex.IMPORT, lex.TYPE:
		return self.parseGlobalWithNoAttr(pub)
	case lex.Attr, lex.FUNC, lex.LET:
		return self.parseGlobalWithAttr(pub)
	default:
		fmt.Println(self.nextTok.Source)
		self.throwErrorf(self.nextTok.Pos, errStrUnknownGlobal)
		return nil
	}
}

// 全局（不带属性）
func (self *Parser) parseGlobalWithNoAttr(pub *lex.Token) Global {
	switch self.nextTok.Kind {
	case lex.IMPORT:
		if pub != nil {
			self.throwErrorf(self.nextTok.Pos, errStrUnknownGlobal)
		}
		return self.parseImport()
	case lex.TYPE:
		return self.parseTypeDef(pub)
	default:
		self.throwErrorf(self.nextTok.Pos, errStrUnknownGlobal)
		return nil
	}
}

// 全局（带属性）
func (self *Parser) parseGlobalWithAttr(pub *lex.Token) Global {
	var attrs []Attr
	if pub == nil {
		for self.nextIs(lex.Attr) {
			if len(attrs) == 0 && pub != nil {
				self.throwErrorf(self.nextTok.Pos, errStrUnknownGlobal)
			}
			attrs = append(attrs, self.parseAttr())
			self.expectNextIs(lex.SEM)
		}

		if self.skipNextIs(lex.PUB) {
			pub = &self.curTok
		}
	}

	switch self.nextTok.Kind {
	case lex.FUNC:
		return self.parseFunction(pub, attrs)
	case lex.LET:
		return self.parseGlobalValue(pub, attrs)
	default:
		self.throwErrorf(self.nextTok.Pos, errStrUnknownGlobal)
		return nil
	}
}

// 包导入
func (self *Parser) parseImport() *Import {
	begin := self.expectNextIs(lex.IMPORT).Pos
	pkgs := self.parseTokenListAtLeastOne(lex.DOT)
	var alias *lex.Token
	if self.skipNextIs(lex.AS) {
		name := self.expectNextIs(lex.IDENT)
		alias = &name
	}
	return NewImport(utils.MixPosition(begin, pkgs[len(pkgs)-1].Pos), pkgs, alias)
}

// 类型定义
func (self *Parser) parseTypeDef(pub *lex.Token) *TypeDef {
	begin := self.expectNextIs(lex.TYPE).Pos
	name := self.expectNextIs(lex.IDENT)
	target := self.parseType()
	if pub == nil {
		return NewTypeDef(utils.MixPosition(begin, target.Position()), false, name, target)
	} else {
		return NewTypeDef(utils.MixPosition(pub.Pos, target.Position()), true, name, target)
	}
}

// 函数
func (self *Parser) parseFunction(pub *lex.Token, attrs []Attr) Global {
	for _, attr := range attrs {
		switch attr.(type) {
		case *AttrExtern, *AttrLink, *AttrNoReturn, *AttrExit, *AttrInline:
		default:
			self.throwErrorf(attr.Position(), errStrCanNotUseAttr)
			return nil
		}
	}

	self.expectNextIs(lex.FUNC)
	var begin utils.Position
	if len(attrs) > 0 {
		begin = attrs[0].Position()
	} else if pub != nil {
		begin = pub.Pos
	} else {
		begin = self.curTok.Pos
	}

	if self.nextIs(lex.LPA) {
		return self.parseMethod(begin, pub != nil, attrs)
	}

	name := self.expectNextIs(lex.IDENT)
	self.expectNextIs(lex.LPA)
	mid := lex.COL
	params := self.parseNameOrNilAndTypeList(&mid, lex.COM, false)
	self.expectNextIs(lex.RPA)
	ret := self.parseTypeOrNil()
	var body *Block
	if self.nextIs(lex.LBR) {
		body = self.parseBlock()
	}
	return NewFunction(utils.MixPosition(begin, self.curTok.Pos), attrs, pub != nil, ret, name, params, body)
}

// 方法
func (self *Parser) parseMethod(begin utils.Position, pub bool, attrs []Attr) *Method {
	for _, attr := range attrs {
		switch attr.(type) {
		case *AttrNoReturn, *AttrExit, *AttrInline:
		default:
			self.throwErrorf(attr.Position(), errStrCanNotUseAttr)
			return nil
		}
	}

	self.expectNextIs(lex.LPA)
	selfTok := self.expectNextIs(lex.IDENT)
	self.expectNextIs(lex.RPA)

	name := self.expectNextIs(lex.IDENT)
	self.expectNextIs(lex.LPA)
	mid := lex.COL
	params := self.parseNameOrNilAndTypeList(&mid, lex.COM, false)
	self.expectNextIs(lex.RPA)
	ret := self.parseTypeOrNil()
	var body *Block
	if self.nextIs(lex.LBR) {
		body = self.parseBlock()
	}
	return NewMethod(utils.MixPosition(begin, self.curTok.Pos), attrs, pub, selfTok, ret, name, params, body)
}

// 全局变量
func (self *Parser) parseGlobalValue(pub *lex.Token, attrs []Attr) *GlobalValue {
	for _, attr := range attrs {
		switch attr.(type) {
		case *AttrExtern, *AttrLink:
		default:
			self.throwErrorf(attr.Position(), errStrCanNotUseAttr)
			return nil
		}
	}
	v := self.parseVariable()
	var begin utils.Position
	if len(attrs) > 0 {
		begin = attrs[0].Position()
	} else if pub != nil {
		begin = pub.Pos
	} else {
		begin = v.Pos
	}
	return NewGlobalValue(utils.MixPosition(begin, v.Position()), attrs, pub != nil, v.Type, v.Name, v.Value)
}
