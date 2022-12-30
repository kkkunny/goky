package parse

import (
	"github.com/kkkunny/klang/src/compiler/lex"
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Attr 属性
type Attr interface {
	Ast
	Attr()
}

// AttrExtern @extern
type AttrExtern struct {
	Pos  utils.Position
	Name lex.Token
}

func NewAttrExtern(pos utils.Position, name lex.Token) *AttrExtern {
	return &AttrExtern{
		Pos:  pos,
		Name: name,
	}
}

func (self AttrExtern) Position() utils.Position {
	return self.Pos
}

func (self AttrExtern) Attr() {}

// AttrLink @link
type AttrLink struct {
	Pos  utils.Position
	Asms []*String
	Libs []*String
}

func NewAttrLink(pos utils.Position, asms, libs []*String) *AttrLink {
	return &AttrLink{
		Pos:  pos,
		Asms: asms,
		Libs: libs,
	}
}

func (self AttrLink) Position() utils.Position {
	return self.Pos
}

func (self AttrLink) Attr() {}

// AttrNoReturn @noreturn
type AttrNoReturn struct {
	Pos utils.Position
}

func NewAttrNoReturn(pos utils.Position) *AttrNoReturn {
	return &AttrNoReturn{Pos: pos}
}

func (self AttrNoReturn) Position() utils.Position {
	return self.Pos
}

func (self AttrNoReturn) Attr() {}

// AttrExit @exit
type AttrExit struct {
	Pos utils.Position
}

func NewAttrExit(pos utils.Position) *AttrExit {
	return &AttrExit{Pos: pos}
}

func (self AttrExit) Position() utils.Position {
	return self.Pos
}

func (self AttrExit) Attr() {}

// ****************************************************************

func (self *Parser) parseAttr() Attr {
	attrName := self.expectNextIs(lex.Attr)
	switch attrName.Source {
	case "@extern":
		self.expectNextIs(lex.LPA)
		name := self.expectNextIs(lex.IDENT)
		end := self.expectNextIs(lex.RPA).Pos
		return NewAttrExtern(utils.MixPosition(attrName.Pos, end), name)
	case "@link":
		self.expectNextIs(lex.LPA)
		var asms, libs []*String
		for {
			linkname := self.expectNextIs(lex.IDENT)
			switch linkname.Source {
			case "asm":
				self.expectNextIs(lex.ASS)
				asms = append(asms, self.parseStringExpr())
			case "lib":
				self.expectNextIs(lex.ASS)
				libs = append(libs, self.parseStringExpr())
			default:
				self.throwErrorf(linkname.Pos, "unknown link")
			}
			if !self.skipNextIs(lex.COM) {
				break
			}
		}
		end := self.expectNextIs(lex.RPA).Pos
		return NewAttrLink(utils.MixPosition(attrName.Pos, end), asms, libs)
	case "@noreturn":
		return NewAttrNoReturn(attrName.Pos)
	case "@exit":
		return NewAttrExit(attrName.Pos)
	default:
		self.throwErrorf(attrName.Pos, "unknown attribute")
		return nil
	}
}
