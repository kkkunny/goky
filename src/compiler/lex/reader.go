package lex

import "io"

// Reader 读取器
type Reader interface {
	io.RuneReader
	io.Seeker
}
