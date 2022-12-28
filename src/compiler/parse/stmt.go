package parse

import (
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Block 代码块
type Block struct {
	Stmts []Stmt `"{" Separator* (@@ (Separator+ @@?)*)? "}"`
}

// Stmt 语句
type Stmt struct {
	utils.Position
	Return      *Return         `@@`
	Variable    *Variable       `| @@`
	Block       *Block          `| @@`
	IfElse      *IfElse         `| @@`
	For         *For            `| @@`
	LoopControl *string         `| @("break" | "continue")`
	Defer       *PrimaryPostfix `| "defer" @@`
	Expr        *Expr           `| @@`
}

// Return 函数返回
type Return struct {
	utils.Position
	Value *Expr `"return" @@?`
}

// Variable 变量
type Variable struct {
	utils.Position
	Name  Name  `"let" @@`
	Type  *Type `(":" @@)?`
	Value *Expr `("=" @@)?`
}

// IfElse 条件分支
type IfElse struct {
	utils.Position
	Cond   Expr          `"if" @@`
	True   Block         `@@`
	Suffix *IfElseSuffix `("else" @@)?`
}

type IfElseSuffix struct {
	False *Block  `@@`
	Else  *IfElse `| @@`
}

// For 循环
type For struct {
	utils.Position
	Cond Expr  `"for" @@`
	Body Block `@@`
}
