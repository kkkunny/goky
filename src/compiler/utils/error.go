package utils

import (
	"fmt"
	"strings"
)

// Error 异常
type Error interface {
	error
	fmt.Stringer
}

// SingleError 单个异常
type SingleError struct {
	Pos Position
	Msg string
}

// Errorf 格式化异常
func Errorf(pos Position, msg string, args ...any) *SingleError {
	return &SingleError{
		Pos: pos,
		Msg: fmt.Sprintf(msg, args...),
	}
}

func (self SingleError) Error() string {
	return self.String()
}

func (self SingleError) String() string {
	return fmt.Sprintf("%s:%d:%d: %s", self.Pos.File, self.Pos.BeginRow, self.Pos.BeginCol, self.Msg)
}

// MultiError 多个异常
type MultiError struct {
	List []Error
}

// NewMultiError 打包异常
func NewMultiError(e ...Error) *MultiError {
	return &MultiError{List: e}
}

func (self MultiError) Error() string {
	return self.String()
}

func (self MultiError) String() string {
	var buf strings.Builder
	for i, e := range self.List {
		buf.WriteString(e.String())
		if i < len(self.List)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
