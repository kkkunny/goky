package parse

import (
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Type 类型
type Type struct {
	utils.Position
	Pointer *Type       `"*" @@`
	Struct  *TypeStruct `| @@`
	Func    *TypeFunc   `| @@`
	Array   *TypeArray  `| @@`
	Tuple   *TypeTuple  `| @@`
	Ident   *TypeIdent  `| @@`
}

// TypeIdent 标识符类型
type TypeIdent struct {
	utils.Position
	Name Name `@@`
}

// TypeFunc 函数类型
type TypeFunc struct {
	utils.Position
	Params TypeList `"func" "(" @@ ")"`
	Ret    *Type    `@@?`
}

// TypeArray 数组类型
type TypeArray struct {
	utils.Position
	Size uint `"[" @Int "]"`
	Elem Type `@@`
}

// TypeTuple 元组类型
type TypeTuple struct {
	utils.Position
	Types TypeList `"(" @@ ")"`
}

// TypeStruct 结构体类型
type TypeStruct struct {
	utils.Position
	Fields []TypeStructField `"struct" "{" Separator* (@@ Separator+)* "}"`
}

type TypeStructField struct {
	Name Name `@@`
	Type Type `":" @@`
}

// TypeList 类型列表
type TypeList struct {
	Types []Type `(@@ ("," @@)*)?`
}
