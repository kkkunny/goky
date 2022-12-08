package parse

import "github.com/kkkunny/klang/src/compiler/internal/utils"

// Block 代码块
type Block struct {
	Stmts []Stmt `"{" Separator* (@@ (Separator* @@)*)? Separator* "}"`
}

// Stmt 语句
type Stmt struct {
	utils.Position
	Return      *Return         `@@ Separator`
	Variable    *Variable       `| @@ Separator`
	Block       *Block          `| @@ Separator`
	IfElse      *IfElse         `| @@ Separator`
	For         *For            `| @@ Separator`
	LoopControl *string         `| @("break" | "continue") Separator`
	Defer       *PrimaryPostfix `| "defer" @@ Separator`
	Expr        *Expr           `| @@ Separator`
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
