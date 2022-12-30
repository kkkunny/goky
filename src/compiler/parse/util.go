package parse

import (
	"bytes"
	"github.com/kkkunny/klang/src/compiler/lex"
	"github.com/kkkunny/stl/list"
	stlos "github.com/kkkunny/stl/os"
	"os"
)

// ParseFile 词法-语法分析单文件
func ParseFile(path stlos.Path) (*File, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(content)
	lexer := lex.NewLexer(path, reader)
	parser := NewParser(lexer)
	return parser.Parse()
}

// ParsePackage 词法-语法分析包
func ParsePackage(path stlos.Path) (*list.SingleLinkedList[*File], error) {
	files, err := os.ReadDir(string(path))
	if err != nil {
		return nil, err
	}
	res := list.NewSingleLinkedList[*File]()
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fp := path.Join(stlos.Path(f.Name()))
		if fp.GetExtension() != "k" {
			continue
		}

		ast, err := ParseFile(fp)
		if err != nil {
			return nil, err
		}
		res.Add(ast)
	}
	return res, nil
}
