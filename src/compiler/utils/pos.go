package utils

import (
	stlos "github.com/kkkunny/stl/os"
)

// Position 位置
type Position struct {
	File             stlos.Path
	Begin, End       uint
	BeginRow, EndRow uint
	BeginCol, EndCol uint
}

// NewPosition 新建Position
func NewPosition(fp stlos.Path) Position {
	return Position{File: fp}
}

// MixPosition 混合Position
func MixPosition(p1, p2 Position) Position {
	if p1.File != p2.File {
		panic("can not mix these positions")
	}
	var b, br, bc uint
	if p1.Begin < p2.Begin {
		b, br, bc = p1.Begin, p1.BeginRow, p1.BeginCol
	} else {
		b, br, bc = p2.Begin, p2.BeginRow, p2.BeginCol
	}
	var e, er, ec uint
	if p1.End > p2.End {
		e, er, ec = p1.End, p1.EndRow, p1.EndCol
	} else {
		e, er, ec = p2.End, p2.EndRow, p2.EndCol
	}
	return Position{
		File:     p1.File,
		Begin:    b,
		End:      e,
		BeginRow: br,
		EndRow:   er,
		BeginCol: bc,
		EndCol:   ec,
	}
}

// SetBegin 设置开始位置
func (self *Position) SetBegin(pos, row, col uint) {
	self.Begin, self.BeginRow, self.BeginCol = pos, row, col
}

// GetBegin 获取开始位置
func (self *Position) GetBegin() (uint, uint, uint) {
	return self.Begin, self.BeginRow, self.BeginCol
}

// SetEnd 设置结束位置
func (self *Position) SetEnd(pos, row, col uint) {
	self.End, self.EndRow, self.EndCol = pos, row, col
}

// GetEnd 获取结束位置
func (self *Position) GetEnd() (uint, uint, uint) {
	return self.End, self.EndRow, self.EndCol
}
