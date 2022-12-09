package parse

import (
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Global 全局
type Global struct {
	GlobalWithAttr *GlobalWithAttr `@@`
	GlobalNoAttr   *GlobalNoAttr   `| @@`
}

// GlobalWithAttr 全局带属性
type GlobalWithAttr struct {
	Attr   []Attr               `(@@ Separator)*`
	Global GlobalWithAttrSuffix `@@`
}

// Attr 属性
type Attr struct {
	utils.Position
	Extern   *Name   `"@extern" "(" @@ ")"`
	LinkAsm  *String `| "@link" "(" "asm" "=" @String ")"`
	LinkLib  *String `| "@link" "(" "lib" "=" @String ")"`
	NoReturn *string `| @"@noreturn"`
	Exit     *string `| @"@exit"`
}

type GlobalWithAttrSuffix struct {
	Function *Function       `@@`
	Variable *GlobalVariable `| @@`
}

// Function 函数
type Function struct {
	utils.Position
	Public *string `@"pub"?`
	Name   Name    `"func" @@`
	Params []Param `"(" (@@ ("," @@)*)? ")"`
	Ret    *Type   `@@?`
	Body   *Block  `@@?`
}

// Param 参数
type Param struct {
	utils.Position
	Name *Name `(@@ ":")?`
	Type Type  `@@`
}

// GlobalVariable 全局变量
type GlobalVariable struct {
	utils.Position
	Public   *string  `@"pub"?`
	Variable Variable `@@`
}

// GlobalNoAttr 全局不带属性
type GlobalNoAttr struct {
	Import *Import `@@`
}

// Import 导入
type Import struct {
	utils.Position
	Packages []Name `"import" @@ ("." @@)*`
	Alias    *Name  `("as" @@)?`
}
