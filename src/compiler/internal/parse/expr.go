package parse

import "github.com/kkkunny/klang/src/compiler/internal/utils"

// Expr 表达式
type Expr struct {
	utils.Position
	Assign Assign `@@`
}

// Assign 赋值
type Assign struct {
	utils.Position
	Left   LogicOpera    `@@`
	Suffix *AssignSuffix `@@?`
}

type AssignSuffix struct {
	Opera string     `@("=" | "+=" | "-=" | "*=" | "/=" | "%=" | "&=" | "|=" | "^=" | "<<=" | ">>=")`
	Right LogicOpera `@@`
}

// LogicOpera 逻辑运算
type LogicOpera struct {
	utils.Position
	Left Equal            `@@`
	Next []LogicOperaNext `@@*`
}

type LogicOperaNext struct {
	utils.Position
	Opera string `@("&&" | "||")`
	Right Equal  `@@`
}

// Equal 加或减
type Equal struct {
	utils.Position
	Left AddOrSub    `@@`
	Next []EqualNext `@@*`
}

type EqualNext struct {
	utils.Position
	Opera string   `@("==" | "!=" | "<" | "<=" | ">" | ">=")`
	Right AddOrSub `@@`
}

// AddOrSub 加或减
type AddOrSub struct {
	utils.Position
	Left MulOrDivOrMod  `@@`
	Next []AddOrSubNext `@@*`
}

type AddOrSubNext struct {
	utils.Position
	Opera string        `@("+" | "-")`
	Right MulOrDivOrMod `@@`
}

// MulOrDivOrMod 乘或除或取余
type MulOrDivOrMod struct {
	utils.Position
	Left ByteOpera           `@@`
	Next []MulOrDivOrModNext `@@*`
}

type MulOrDivOrModNext struct {
	utils.Position
	Opera string    `@("*" | "/" | "%")`
	Right ByteOpera `@@`
}

// ByteOpera 位运算
type ByteOpera struct {
	utils.Position
	Left UnaryPostfix    `@@`
	Next []ByteOperaNext `@@*`
}

type ByteOperaNext struct {
	utils.Position
	Opera string       `@("&" | "|" | "^" | "<<" | ">>")`
	Right UnaryPostfix `@@`
}

// UnaryPostfix 后缀
type UnaryPostfix struct {
	utils.Position
	Unary  Unary                `@@`
	Suffix []UnaryPostfixSuffix `@@*`
}

type UnaryPostfixSuffix struct {
	As     *Type                     `"as" @@`
	Select *UnaryPostfixSuffixSelect `| "?" @@`
}

type UnaryPostfixSuffixSelect struct {
	True  Expr         `@@`
	False UnaryPostfix `":" @@`
}

// Unary 一元运算
type Unary struct {
	utils.Position
	Opera   *string        `@("-" | "~" | "!" | "&" | "*")?`
	Postfix PrimaryPostfix `@@`
}

// PrimaryPostfix 后缀
type PrimaryPostfix struct {
	utils.Position
	Primary Primary                `@@`
	Suffix  []PrimaryPostfixSuffix `@@*`
}

type PrimaryPostfixSuffix struct {
	Call  *ExprList `"(" @@ ")"`
	Index *Expr     `| "[" @@ "]"`
	Dot   *Name     `| "." @@`
}

// Primary 单表达式
type Primary struct {
	utils.Position
	Constant *Constant `@@`
	Ident    *Ident    `| @@`
	Tuple    *ExprList `| "(" @@ ")"`
	Array    *ExprList `| "[" @@ "]"`
	Struct   *ExprList `| "{" @@ "}"`
}

// Constant 常量
type Constant struct {
	utils.Position
	Null    *string  `@"null"`
	Int     *int64   `| @Int`
	Float   *float64 `| @Float`
	Bool    *Bool    `| @("true" | "false")`
	Char    *Char    `| @Char`
	CString *CString `| @CString`
	String  *String  `| @String`
}

// Ident 标识符
type Ident struct {
	utils.Position
	Package *Name `(@@ ":")?`
	Name    Name  `@@`
}

// ExprList 表达式列表
type ExprList struct {
	Exprs []Expr `(@@ ("," @@)*)?`
}
