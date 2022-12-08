package parse

import (
	"github.com/kkkunny/klang/src/compiler/internal/utils"
	"strings"
)

type Name struct {
	utils.Position
	Value string `@Name`
}

type Bool bool

func (self *Bool) Capture(values []string) error {
	*self = values[0] == "true"
	return nil
}

type Char rune

func (self *Char) Capture(values []string) error {
	s := values[0][1 : len(values[0])-1]
	switch s {
	case "\\0":
		*self = 0
	case "\\a":
		*self = 7
	case "\\b":
		*self = 8
	case "\\t":
		*self = 9
	case "\\n":
		*self = 10
	case "\\v":
		*self = 11
	case "\\f":
		*self = 12
	case "\\r":
		*self = 13
	case "\\\\":
		*self = '\\'
	case "\\'":
		*self = '\''
	default:
		*self = Char([]rune(s)[0])
	}
	return nil
}

type String []rune

func (self *String) Capture(values []string) error {
	s := values[0][1 : len(values[0])-1]
	r := strings.NewReplacer(
		"\\0", string([]byte{0}),
		"\\a", string([]byte{7}),
		"\\b", string([]byte{8}),
		"\\t", string([]byte{9}),
		"\\n", string([]byte{10}),
		"\\v", string([]byte{11}),
		"\\f", string([]byte{12}),
		"\\r", string([]byte{13}),
		"\\\\", "\\",
		"\\\"", "\"",
	)
	*self = String(r.Replace(s))
	return nil
}

type CString []byte

func (self *CString) Capture(values []string) error {
	s := values[0][2 : len(values[0])-1]
	r := strings.NewReplacer(
		"\\0", string([]byte{0}),
		"\\a", string([]byte{7}),
		"\\b", string([]byte{8}),
		"\\t", string([]byte{9}),
		"\\n", string([]byte{10}),
		"\\v", string([]byte{11}),
		"\\f", string([]byte{12}),
		"\\r", string([]byte{13}),
		"\\\\", "\\",
		"\\\"", "\"",
	)
	*self = CString(r.Replace(s))
	*self = append(*self, 0)
	return nil
}
