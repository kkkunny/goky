package parse

import (
	"bytes"
	"github.com/kkkunny/klang/src/compiler/lex"
	stlos "github.com/kkkunny/stl/os"
	"os"
)

// ParseFile 词法-语法分析单文件
func ParseFile(path stlos.Path) (*Package, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(content)
	lexer := lex.NewLexer(path, reader)
	parser := NewParser(lexer)
	file, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return NewPackage(file.Path.GetParent(), file), nil
}

// ParsePackage 词法-语法分析包
func ParsePackage(path stlos.Path) (*Package, error) {
	files, err := os.ReadDir(string(path))
	if err != nil {
		return nil, err
	}
	pkg := NewPackage(path)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fp := path.Join(stlos.Path(f.Name()))
		if fp.GetExtension() != "k" {
			continue
		}

		file, err := ParseFile(fp)
		if err != nil {
			return nil, err
		}
		pkg.Files = append(pkg.Files, file.Files[0])
	}
	return pkg, nil
}
