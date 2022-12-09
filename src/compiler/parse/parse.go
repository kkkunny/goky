package parse

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/kkkunny/klang/src/compiler/utils"
	stlos "github.com/kkkunny/stl/os"
	"os"
)

var (
	lmg = lexer.MustSimple([]lexer.SimpleRule{
		{"Whitespace", `[ \t]`},
		{"Comment", `//.*\n`},
		{"Float", `[0-9].[0-9]+`},
		{"Int", `[0-9]+`},
		{"Char", `'\\''|'\\0'|'\\a'|'\\b'|'\\t'|'\\n'|'\\v'|'\\f'|'\\r'|'\\\\'|'.'`},
		{"CString", `c".*"`},
		{"String", `".*"`},
		{"Separator", `[;\n]`},
		{"Attr", `@[a-zA-Z_]+`},
		{"Name", `[a-zA-Z_][a-zA-Z0-9_]*`},
		{"Symbol", `::|\+=|-=|\*=|/=|%=|&=|\|=|\^=|<<=|>>=|&&|\|\||<=|>=|==|!=|<=|>=|<<|>>|\(|\)|\{|\}|\+|-|\*|/|%|&|\||\^|~|:|=|,|\[|\]|\.|<|>|!|\?`},
	})
	pmg = participle.MustBuild[Package](
		participle.Lexer(lmg),
		participle.Elide("Whitespace", "Comment"),
		participle.UseLookahead(2),
	)
)

// Package 包
type Package struct {
	PkgPath stlos.Path
	Globals []Global `Separator* (@@ (Separator+ @@)*)? Separator*`
}

// ParseFile 词法分析+语法分析 文件
func ParseFile(fp stlos.Path) (*Package, error) {
	fp, err := fp.GetAbsolute()
	if err != nil {
		return nil, err
	}
	if !fp.IsFile() || fp.GetExtension() != "k" {
		return nil, fmt.Errorf("expect a k source file")
	}
	file, err := os.Open(fp.String())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	pkg, err := pmg.Parse(fp.String(), file)
	if err != nil {
		if parseErr, ok := err.(participle.Error); ok {
			return nil, utils.Errorf(utils.NewPosition(parseErr.Position()), parseErr.Message())
		} else {
			panic(err)
		}
	}
	pkg.PkgPath = fp.GetParent()
	return pkg, nil
}

// ParsePackage 词法分析+语法分析 包
func ParsePackage(fp stlos.Path) (*Package, error) {
	fp, err := fp.GetAbsolute()
	if err != nil {
		return nil, err
	}
	if !fp.IsDir() {
		return nil, fmt.Errorf("expect a goky source package")
	}
	fileInfoes, err := os.ReadDir(fp.String())
	if err != nil {
		return nil, err
	}
	var globals []Global
	for _, fileInfo := range fileInfoes {
		filepath := fp.Join(stlos.Path(fileInfo.Name()))
		if filepath.GetExtension() != "k" {
			continue
		}
		pkg, err := ParseFile(filepath)
		if err != nil {
			return nil, err
		}
		globals = append(globals, pkg.Globals...)
	}
	return &Package{
		PkgPath: fp,
		Globals: globals,
	}, nil
}
