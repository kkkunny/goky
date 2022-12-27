package analyse

import (
	"github.com/kkkunny/klang/src/compiler/parse"
	"github.com/kkkunny/klang/src/compiler/utils"
)

// Expr 表达式
type Expr interface {
	Stmt
	GetType() Type
	GetMut() bool
	IsTemporary() bool
}

// Ident 标识符
type Ident interface {
	Expr
	ident()
}

// Integer 整数
type Integer struct {
	Type  Type
	Value int64
}

func (self Integer) stmt() {}

func (self Integer) GetType() Type {
	return self.Type
}

func (self Integer) GetMut() bool {
	return false
}

func (self Integer) IsTemporary() bool {
	return true
}

// Float 浮点数
type Float struct {
	Type  Type
	Value float64
}

func (self Float) stmt() {}

func (self Float) GetType() Type {
	return self.Type
}

func (self Float) GetMut() bool {
	return false
}

func (self Float) IsTemporary() bool {
	return true
}

// Boolean 布尔数
type Boolean struct {
	Type  Type
	Value bool
}

func (self Boolean) stmt() {}

func (self Boolean) GetType() Type {
	return self.Type
}

func (self Boolean) GetMut() bool {
	return false
}

func (self Boolean) IsTemporary() bool {
	return true
}

// Null 空指针
type Null struct {
	Type Type
}

func (self Null) stmt() {}

func (self Null) GetType() Type {
	return self.Type
}

func (self Null) GetMut() bool {
	return false
}

func (self Null) IsTemporary() bool {
	return true
}

// Binary 二元表达式
type Binary struct {
	Opera       string
	Left, Right Expr
}

func (self Binary) stmt() {}

func (self Binary) GetType() Type {
	return self.Left.GetType()
}

func (self Binary) GetMut() bool {
	return false
}

func (self Binary) IsTemporary() bool {
	return true
}

// FuncCall 函数调用
type FuncCall struct {
	NoReturn bool
	Exit     bool

	Func Expr
	Args []Expr
}

func (self FuncCall) stmt() {}

func (self FuncCall) GetType() Type {
	return GetBaseType(self.Func.GetType()).(*TypeFunc).Ret
}

func (self FuncCall) GetMut() bool {
	return false
}

func (self FuncCall) IsTemporary() bool {
	return true
}

// MethodCall 方法调用
type MethodCall struct {
	NoReturn bool
	Exit     bool

	Method *Method
	Args   []Expr
}

func (self MethodCall) stmt() {}

func (self MethodCall) GetType() Type {
	return self.Method.GetType().(*TypeFunc).Ret
}

func (self MethodCall) GetMut() bool {
	return false
}

func (self MethodCall) IsTemporary() bool {
	return true
}

// Param 参数
type Param struct {
	Type Type
}

func (self Param) stmt() {}

func (self Param) ident() {}

func (self Param) GetType() Type {
	return self.Type
}

func (self Param) GetMut() bool {
	return true
}

func (self Param) IsTemporary() bool {
	return false
}

// Array 数组
type Array struct {
	Type  Type
	Elems []Expr
}

func (self Array) stmt() {}

func (self Array) GetType() Type {
	return self.Type
}

func (self Array) GetMut() bool {
	return false
}

func (self Array) IsTemporary() bool {
	return true
}

// EmptyArray 空数组
type EmptyArray struct {
	Type Type
}

func (self EmptyArray) stmt() {}

func (self EmptyArray) GetType() Type {
	return self.Type
}

func (self EmptyArray) GetMut() bool {
	return false
}

func (self EmptyArray) IsTemporary() bool {
	return true
}

// Assign 赋值
type Assign struct {
	Opera       string
	Left, Right Expr
}

func (self Assign) stmt() {}

func (self Assign) GetType() Type {
	return None
}

func (self Assign) GetMut() bool {
	return false
}

func (self Assign) IsTemporary() bool {
	return true
}

// Equal 赋值
type Equal struct {
	Opera       string
	Left, Right Expr
}

func (self Equal) stmt() {}

func (self Equal) GetType() Type {
	return Bool
}

func (self Equal) GetMut() bool {
	return false
}

func (self Equal) IsTemporary() bool {
	return true
}

// Unary 一元表达式
type Unary struct {
	Type  Type
	Opera string
	Value Expr
}

func (self Unary) stmt() {}

func (self Unary) GetType() Type {
	return self.Type
}

func (self Unary) GetMut() bool {
	return false
}

func (self Unary) IsTemporary() bool {
	return true
}

// Index 索引
type Index struct {
	Type        Type
	From, Index Expr
}

func (self Index) stmt() {}

func (self Index) GetType() Type {
	return self.Type
}

func (self Index) GetMut() bool {
	return self.From.GetMut()
}

func (self Index) IsTemporary() bool {
	return self.From.IsTemporary()
}

// Select 选择
type Select struct {
	Cond, True, False Expr
}

func (self Select) stmt() {}

func (self Select) GetType() Type {
	return self.True.GetType()
}

func (self Select) GetMut() bool {
	return self.True.GetMut() && self.False.GetMut()
}

func (self Select) IsTemporary() bool {
	return self.True.IsTemporary() || self.False.IsTemporary()
}

// Tuple 元组
type Tuple struct {
	Type  Type
	Elems []Expr
}

func (self Tuple) stmt() {}

func (self Tuple) GetType() Type {
	return self.Type
}

func (self Tuple) GetMut() bool {
	return false
}

func (self Tuple) IsTemporary() bool {
	return true
}

// EmptyTuple 空元组
type EmptyTuple struct {
	Type Type
}

func (self EmptyTuple) stmt() {}

func (self EmptyTuple) GetType() Type {
	return self.Type
}

func (self EmptyTuple) GetMut() bool {
	return false
}

func (self EmptyTuple) IsTemporary() bool {
	return true
}

// Struct 结构体
type Struct struct {
	Type   Type
	Fields []Expr
}

func (self Struct) stmt() {}

func (self Struct) GetType() Type {
	return self.Type
}

func (self Struct) GetMut() bool {
	return false
}

func (self Struct) IsTemporary() bool {
	return true
}

// EmptyStruct 空结构体
type EmptyStruct struct {
	Type Type
}

func (self EmptyStruct) stmt() {}

func (self EmptyStruct) GetType() Type {
	return self.Type
}

func (self EmptyStruct) GetMut() bool {
	return false
}

func (self EmptyStruct) IsTemporary() bool {
	return true
}

// GetField 获取成员
type GetField struct {
	From  Expr
	Index string
}

func (self GetField) stmt() {}

func (self GetField) GetType() Type {
	return GetBaseType(self.From.GetType()).(*TypeStruct).Fields.Get(self.Index)
}

func (self GetField) GetMut() bool {
	return self.From.GetMut()
}

func (self GetField) IsTemporary() bool {
	return self.From.IsTemporary()
}

// Covert 类型转换
type Covert struct {
	From Expr
	To   Type
}

func (self Covert) stmt() {}

func (self Covert) GetType() Type {
	return self.To
}

func (self Covert) GetMut() bool {
	return false
}

func (self Covert) IsTemporary() bool {
	return true
}

// Method 方法
type Method struct {
	Self Expr // 类型定义 || 类型定义指针
	Func *Function
}

func (self Method) stmt() {}

func (self Method) GetType() Type {
	return self.Func.GetType()
}

func (self Method) GetMut() bool {
	return false
}

func (self Method) IsTemporary() bool {
	return true
}

// *********************************************************************************************************************

// 表达式
func analyseExpr(ctx *blockContext, expect Type, ast parse.Expr) (Expr, utils.Error) {
	return analyseAssign(ctx, expect, ast.Assign)
}

// 赋值
func analyseAssign(ctx *blockContext, expect Type, ast parse.Assign) (Expr, utils.Error) {
	if ast.Suffix == nil {
		return analyseLogicOpera(ctx, expect, ast.Left)
	} else {
		left, err := analyseLogicOpera(ctx, expect, ast.Left)
		if err != nil {
			return nil, err
		} else if !left.GetMut() {
			return nil, utils.Errorf(ast.Left.Position, "expect a mutable value")
		}
		lt := left.GetType()
		switch ast.Suffix.Opera {
		case "+=", "-=", "*=", "/=", "%=":
			if !IsNumberTypeAndSon(lt) {
				return nil, utils.Errorf(ast.Left.Position, "expect a number")
			}
		case "&=", "|=", "^=", "<<=", ">>=":
			if !IsIntTypeAndSon(lt) {
				return nil, utils.Errorf(ast.Left.Position, "expect a integer")
			}
		}
		right, err := analyseLogicOpera(ctx, lt, ast.Suffix.Right)
		if err != nil {
			return nil, err
		}
		right, err = expectExpr(ast.Suffix.Right.Position, lt, right)
		if err != nil {
			return nil, err
		}
		return &Assign{
			Opera: ast.Suffix.Opera,
			Left:  left,
			Right: right,
		}, nil
	}
}

// 逻辑运算
func analyseLogicOpera(ctx *blockContext, expect Type, ast parse.LogicOpera) (Expr, utils.Error) {
	if len(ast.Next) == 0 {
		return analyseEqual(ctx, expect, ast.Left)
	} else {
		leftPos := ast.Position
		left, err := analyseEqual(ctx, expect, ast.Left)
		if err != nil {
			return nil, err
		}
		left, err = expectExprAndSon(ast.Left.Position, Bool, left)
		if err != nil {
			return nil, err
		}
		for _, next := range ast.Next {
			right, err := analyseEqual(ctx, left.GetType(), next.Right)
			if err != nil {
				return nil, err
			}
			right, err = expectExpr(next.Right.Position, left.GetType(), right)
			if err != nil {
				return nil, err
			}
			left = &Binary{
				Opera: next.Opera,
				Left:  left,
				Right: right,
			}
			leftPos = utils.MixPosition(leftPos, next.Position)
		}
		return left, nil
	}
}

// 比较
func analyseEqual(ctx *blockContext, expect Type, ast parse.Equal) (Expr, utils.Error) {
	if len(ast.Next) == 0 {
		return analyseAddOrSub(ctx, expect, ast.Left)
	} else {
		leftPos := ast.Position
		left, err := analyseAddOrSub(ctx, nil, ast.Left)
		if err != nil {
			return nil, err
		}
		for _, next := range ast.Next {
			lt := left.GetType()
			if IsNoneType(lt) {
				return nil, utils.Errorf(leftPos, "expect a value")
			} else if next.Opera != "==" && next.Opera != "!=" && !IsNumberTypeAndSon(lt) {
				return nil, utils.Errorf(leftPos, "expect a number")
			}
			right, err := analyseAddOrSub(ctx, lt, next.Right)
			if err != nil {
				return nil, err
			}
			right, err = expectExpr(next.Right.Position, lt, right)
			if err != nil {
				return nil, err
			}
			left = &Equal{
				Opera: next.Opera,
				Left:  left,
				Right: right,
			}
			leftPos = utils.MixPosition(leftPos, next.Position)
		}
		return left, nil
	}
}

// 加或减
func analyseAddOrSub(ctx *blockContext, expect Type, ast parse.AddOrSub) (Expr, utils.Error) {
	if len(ast.Next) == 0 {
		return analyseMulOrDivOrMod(ctx, expect, ast.Left)
	} else {
		leftPos := ast.Position
		left, err := analyseMulOrDivOrMod(ctx, expect, ast.Left)
		if err != nil {
			return nil, err
		}
		for _, next := range ast.Next {
			lt := left.GetType()
			if !IsNumberTypeAndSon(lt) {
				return nil, utils.Errorf(leftPos, "expect a number")
			}
			right, err := analyseMulOrDivOrMod(ctx, lt, next.Right)
			if err != nil {
				return nil, err
			}
			right, err = expectExpr(next.Right.Position, lt, right)
			if err != nil {
				return nil, err
			}
			left = &Binary{
				Opera: next.Opera,
				Left:  left,
				Right: right,
			}
			leftPos = utils.MixPosition(leftPos, next.Position)
		}
		return left, nil
	}
}

// 乘或除或取余
func analyseMulOrDivOrMod(ctx *blockContext, expect Type, ast parse.MulOrDivOrMod) (Expr, utils.Error) {
	if len(ast.Next) == 0 {
		return analyseByteOpera(ctx, expect, ast.Left)
	} else {
		leftPos := ast.Position
		left, err := analyseByteOpera(ctx, expect, ast.Left)
		if err != nil {
			return nil, err
		}
		for _, next := range ast.Next {
			lt := left.GetType()
			if !IsNumberTypeAndSon(lt) {
				return nil, utils.Errorf(leftPos, "expect a number")
			}
			right, err := analyseByteOpera(ctx, lt, next.Right)
			if err != nil {
				return nil, err
			}
			right, err = expectExpr(next.Right.Position, lt, right)
			if err != nil {
				return nil, err
			}
			left = &Binary{
				Opera: next.Opera,
				Left:  left,
				Right: right,
			}
			leftPos = utils.MixPosition(leftPos, next.Position)
		}
		return left, nil
	}
}

// 位运算
func analyseByteOpera(ctx *blockContext, expect Type, ast parse.ByteOpera) (Expr, utils.Error) {
	if len(ast.Next) == 0 {
		return analyseUnaryPostfix(ctx, expect, ast.Left)
	} else {
		leftPos := ast.Position
		left, err := analyseUnaryPostfix(ctx, expect, ast.Left)
		if err != nil {
			return nil, err
		}
		for _, next := range ast.Next {
			lt := left.GetType()
			if !IsIntTypeAndSon(lt) {
				return nil, utils.Errorf(leftPos, "expect a integer")
			}
			right, err := analyseUnaryPostfix(ctx, lt, next.Right)
			if err != nil {
				return nil, err
			}
			right, err = expectExpr(next.Right.Position, lt, right)
			if err != nil {
				return nil, err
			}
			left = &Binary{
				Opera: next.Opera,
				Left:  left,
				Right: right,
			}
			leftPos = utils.MixPosition(leftPos, next.Position)
		}
		return left, nil
	}
}

// 一元后缀
func analyseUnaryPostfix(ctx *blockContext, expect Type, ast parse.UnaryPostfix) (Expr, utils.Error) {
	if ast.Suffix == nil {
		return analyseUnary(ctx, expect, ast.Unary)
	} else {
		var left Expr
		var leftPos utils.Position
		for _, suffix := range ast.Suffix {
			switch {
			case suffix.Select != nil:
				var err utils.Error
				if left == nil {
					left, err = analyseUnary(ctx, Bool, ast.Unary)
					if err != nil {
						return nil, err
					}
					leftPos = ast.Unary.Position
				}
				cond, err := expectExprAndSon(leftPos, Bool, left)
				if err != nil {
					return nil, err
				}
				tv, err := analyseExpr(ctx, expect, suffix.Select.True)
				if err != nil {
					return nil, err
				}
				fv, err := analyseUnaryPostfix(ctx, tv.GetType(), suffix.Select.False)
				if err != nil {
					return nil, err
				}
				fv, err = expectExpr(suffix.Select.False.Position, tv.GetType(), fv)
				if err != nil {
					return nil, err
				}
				left = &Select{
					Cond:  cond,
					True:  tv,
					False: fv,
				}
				leftPos = utils.MixPosition(leftPos, suffix.Select.False.Position)
			case suffix.As != nil:
				to, err := analyseType(ctx.GetPackageContext(), suffix.As)
				if err != nil {
					return nil, err
				}
				if left == nil {
					left, err = analyseUnary(ctx, to, ast.Unary)
					if err != nil {
						return nil, err
					}
					leftPos = ast.Unary.Position
				}
				ft := left.GetType()

				switch {
				case GetDepthBaseType(ft).Equal(GetDepthBaseType(to)):
				case IsNumberTypeAndSon(ft) && IsNumberTypeAndSon(to):
				case GetBaseType(ft).Equal(Usize) && (IsPtrTypeAndSon(to) || IsFuncTypeAndSon(to)):
				case (IsPtrTypeAndSon(ft) || IsFuncTypeAndSon(ft)) && GetBaseType(to).Equal(Usize):
				case (IsPtrTypeAndSon(ft) || IsFuncTypeAndSon(ft)) && (IsPtrTypeAndSon(to) || IsFuncTypeAndSon(to)):
				default:
					return nil, utils.Errorf(leftPos, "can not covert to type `%s`", to)
				}

				left = &Covert{
					From: left,
					To:   to,
				}
				leftPos = utils.MixPosition(leftPos, suffix.As.Position)
			default:
				panic("")
			}
		}
		return left, nil
	}
}

// 一元运算
func analyseUnary(ctx *blockContext, expect Type, ast parse.Unary) (Expr, utils.Error) {
	if ast.Opera == nil {
		return analysePrimaryPostfix(ctx, expect, ast.Postfix)
	} else {
		switch *ast.Opera {
		case "-":
			value, err := analysePrimaryPostfix(ctx, expect, ast.Postfix)
			if err != nil {
				return nil, err
			}
			if !IsNumberTypeAndSon(value.GetType()) {
				return nil, utils.Errorf(ast.Postfix.Position, "expect a number")
			}
			return &Binary{
				Opera: "-",
				Left:  getDefaultExprByType(value.GetType()),
				Right: value,
			}, nil
		case "~":
			value, err := analysePrimaryPostfix(ctx, expect, ast.Postfix)
			if err != nil {
				return nil, err
			}
			if !IsSintTypeAndSon(value.GetType()) {
				return nil, utils.Errorf(ast.Postfix.Position, "expect a signed integer")
			}
			return &Binary{
				Opera: "^",
				Left:  value,
				Right: &Integer{
					Type:  value.GetType(),
					Value: -1,
				},
			}, nil
		case "!":
			value, err := analysePrimaryPostfix(ctx, expect, ast.Postfix)
			if err != nil {
				return nil, err
			}
			value, err = expectExprAndSon(ast.Postfix.Position, Bool, value)
			if err != nil {
				return nil, err
			}
			return &Unary{
				Type:  value.GetType(),
				Opera: "!",
				Value: value,
			}, nil
		case "&":
			if expect != nil && IsPtrTypeAndSon(expect) {
				expect = GetBaseType(expect).(*TypePtr).Elem
			}
			value, err := analysePrimaryPostfix(ctx, expect, ast.Postfix)
			if err != nil {
				return nil, err
			}
			if value.IsTemporary() {
				return nil, utils.Errorf(ast.Postfix.Position, "not expect a temporary value")
			}
			return &Unary{
				Type:  NewPtrType(value.GetType()),
				Opera: "&",
				Value: value,
			}, nil
		case "*":
			if expect != nil {
				expect = NewPtrType(expect)
			}
			value, err := analysePrimaryPostfix(ctx, expect, ast.Postfix)
			if err != nil {
				return nil, err
			}
			vt := value.GetType()
			if !IsPtrTypeAndSon(vt) {
				return nil, utils.Errorf(ast.Postfix.Position, "expect a pointer")
			}
			return &Unary{
				Type:  GetBaseType(vt).(*TypePtr).Elem,
				Opera: "*",
				Value: value,
			}, nil
		default:
			panic("")
		}
	}
}

// 单表达式后缀
func analysePrimaryPostfix(ctx *blockContext, expect Type, ast parse.PrimaryPostfix) (Expr, utils.Error) {
	if ast.Suffix == nil {
		return analysePrimary(ctx, expect, ast.Primary)
	} else if len(ast.Suffix) == 1 {
		return analyseSinglePrimaryPostfix(ctx, expect, parse.PrimaryPostfix{
			Position: ast.Primary.Position,
			Primary:  ast.Primary,
		}, ast.Suffix[0])
	} else {
		prefix := parse.PrimaryPostfix{
			Position: ast.Primary.Position,
			Primary:  ast.Primary,
		}
		for _, suffix := range ast.Suffix[:len(ast.Suffix)-1] {
			prefix = parse.PrimaryPostfix{
				Position: prefix.Position,
				Primary: parse.Primary{
					Position: prefix.Position,
					Tuple: &parse.ExprList{Exprs: []parse.Expr{
						{
							Position: prefix.Position,
							Assign: parse.Assign{
								Position: prefix.Position,
								Left: parse.LogicOpera{
									Left: parse.Equal{
										Position: prefix.Position,
										Left: parse.AddOrSub{
											Position: prefix.Position,
											Left: parse.MulOrDivOrMod{
												Position: prefix.Position,
												Left: parse.ByteOpera{
													Position: prefix.Position,
													Left: parse.UnaryPostfix{
														Unary: parse.Unary{
															Position: prefix.Position,
															Postfix:  prefix,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}},
				},
				Suffix: []parse.PrimaryPostfixSuffix{suffix},
			}
		}
		return analyseSinglePrimaryPostfix(ctx, expect, prefix, ast.Suffix[len(ast.Suffix)-1])
	}
}

// 单个单表达式后缀
func analyseSinglePrimaryPostfix(ctx *blockContext, expect Type, prefixAst parse.PrimaryPostfix, suffixAst parse.PrimaryPostfixSuffix) (Expr, utils.Error) {
	switch {
	case suffixAst.Call != nil:
		f, err := analysePrimaryPostfix(ctx, nil, prefixAst)
		if err != nil {
			return nil, err
		}
		ft, ok := GetBaseType(f.GetType()).(*TypeFunc)
		if !ok {
			return nil, utils.Errorf(prefixAst.Position, "expect a function")
		}

		if method, ok := f.(*Method); ok {
			if len(ft.Params)-1 != len(suffixAst.Call.Exprs) {
				return nil, utils.Errorf(prefixAst.Position, "expect %d arguments", len(ft.Params)-1)
			}
			args, err := analyseExprList(ctx, ft.Params[1:], *suffixAst.Call)
			if err != nil {
				return nil, err
			}
			var errors []utils.Error
			for i, pt := range ft.Params[1:] {
				var err utils.Error
				args[i], err = expectExpr(suffixAst.Call.Exprs[i].Position, pt, args[i])
				if err != nil {
					errors = append(errors, err)
				}
			}
			if len(errors) == 1 {
				return nil, errors[0]
			} else if len(errors) > 1 {
				return nil, utils.NewMultiError(errors...)
			}

			var noReturn bool
			var exit bool
			if g, ok := f.(*Function); ok {
				if g.Exit {
					exit = true
				}
				if g.NoReturn {
					noReturn = true
					ctx.SetEnd()
				}
			}

			return &MethodCall{
				NoReturn: noReturn,
				Exit:     exit,
				Method:   method,
				Args:     args,
			}, nil
		} else {
			if len(ft.Params) != len(suffixAst.Call.Exprs) {
				return nil, utils.Errorf(prefixAst.Position, "expect %d arguments", len(ft.Params))
			}
			args, err := analyseExprList(ctx, ft.Params, *suffixAst.Call)
			if err != nil {
				return nil, err
			}
			var errors []utils.Error
			for i, a := range args {
				a, err = expectExpr(suffixAst.Call.Exprs[i].Position, ft.Params[i], a)
				if err != nil {
					errors = append(errors, err)
				}
			}
			if len(errors) == 1 {
				return nil, errors[0]
			} else if len(errors) > 1 {
				return nil, utils.NewMultiError(errors...)
			}

			var noReturn bool
			var exit bool
			if g, ok := f.(*Function); ok {
				if g.Exit {
					exit = true
				}
				if g.NoReturn {
					noReturn = true
					ctx.SetEnd()
				}
			}

			return &FuncCall{
				NoReturn: noReturn,
				Exit:     exit,
				Func:     f,
				Args:     args,
			}, nil
		}
	case suffixAst.Index != nil:
		// TODO expect
		prefix, err := analysePrimaryPostfix(ctx, nil, prefixAst)
		if err != nil {
			return nil, err
		}
		switch pt := GetBaseType(prefix.GetType()).(type) {
		case *TypeArray:
			index, err := analyseExpr(ctx, Usize, *suffixAst.Index)
			if err != nil {
				return nil, err
			}
			index, err = expectExprAndSon(suffixAst.Index.Position, Usize, index)
			if err != nil {
				return nil, err
			}
			return &Index{
				Type:  pt.Elem,
				From:  prefix,
				Index: index,
			}, nil
		case *TypePtr:
			index, err := analyseExpr(ctx, Usize, *suffixAst.Index)
			if err != nil {
				return nil, err
			}
			index, err = expectExprAndSon(suffixAst.Index.Position, Usize, index)
			if err != nil {
				return nil, err
			}
			return &Index{
				Type:  pt.Elem,
				From:  prefix,
				Index: index,
			}, nil
		case *TypeTuple:
			index, err := analyseExpr(ctx, Usize, *suffixAst.Index)
			if err != nil {
				return nil, err
			}
			literal, ok := index.(*Integer)
			if !ok {
				return nil, utils.Errorf(suffixAst.Index.Position, "expect a integer literal")
			}
			return &Index{
				Type:  pt.Elems[literal.Value],
				From:  prefix,
				Index: literal,
			}, nil
		default:
			return nil, utils.Errorf(prefixAst.Position, "expect a array or tuple")
		}
	case suffixAst.Dot != nil:
		prefix, err := analysePrimaryPostfix(ctx, nil, prefixAst)
		if err != nil {
			return nil, err
		}

		// 方法
		prefixType := prefix.GetType()
		if IsTypedef(prefixType) || (IsPtrType(prefixType) && IsTypedef(prefixType.(*TypePtr).Elem)) {
			var _selfType *Typedef
			if td, ok := prefixType.(*Typedef); ok {
				_selfType = td
			} else {
				_selfType = prefixType.(*TypePtr).Elem.(*Typedef)
			}

			funcName := _selfType.String() + "." + suffixAst.Dot.Value
			if funcObj := ctx.GetValue(funcName); funcObj != nil {
				fun := funcObj.(*Function)
				return &Method{
					Self: prefix,
					Func: fun,
				}, nil
			}
		}

		// 属性
		switch t := GetBaseType(prefixType).(type) {
		case *TypeStruct:
			if !t.Fields.ContainKey(suffixAst.Dot.Value) {
				return nil, utils.Errorf(suffixAst.Dot.Position, "unknown identifier")
			}
			return &GetField{
				From:  prefix,
				Index: suffixAst.Dot.Value,
			}, nil
		case *TypePtr:
			st, ok := GetBaseType(t.Elem).(*TypeStruct)
			if !ok {
				break
			}
			if !st.Fields.ContainKey(suffixAst.Dot.Value) {
				return nil, utils.Errorf(suffixAst.Dot.Position, "unknown identifier")
			}
			return &GetField{
				From: &Unary{
					Type:  t.Elem,
					Opera: "*",
					Value: prefix,
				},
				Index: suffixAst.Dot.Value,
			}, nil
		}
		return nil, utils.Errorf(prefixAst.Position, "expect a struct")
	default:
		panic("")
	}
}

// 单表达式
func analysePrimary(ctx *blockContext, expect Type, ast parse.Primary) (Expr, utils.Error) {
	switch {
	case ast.Constant != nil:
		return analyseConstant(expect, *ast.Constant)
	case ast.Ident != nil:
		return analyseIdent(ctx, *ast.Ident)
	case ast.Tuple != nil:
		if len(ast.Tuple.Exprs) == 0 {
			if expect == nil || !IsTupleTypeAndSon(expect) {
				return nil, utils.Errorf(ast.Position, "expect a tuple type")
			}
			return &EmptyTuple{Type: expect}, nil
		} else if len(ast.Tuple.Exprs) == 1 && (expect == nil || !IsTupleTypeAndSon(expect) || len(GetBaseType(expect).(*TypeTuple).Elems) != 1) {
			return analyseExpr(ctx, expect, ast.Tuple.Exprs[0])
		}
		expects := make([]Type, len(ast.Tuple.Exprs))
		if expect != nil {
			if tt, ok := GetBaseType(expect).(*TypeTuple); ok && len(tt.Elems) == len(ast.Tuple.Exprs) {
				for i := range expects {
					expects[i] = tt.Elems[i]
				}
			}
		}
		elems, err := analyseExprList(ctx, expects, *ast.Tuple)
		if err != nil {
			return nil, err
		}
		for i, e := range elems {
			expects[i] = e.GetType()
		}
		var rt Type = NewTupleType(expects...)
		if expect != nil && GetDepthBaseType(expect).Equal(GetDepthBaseType(rt)) {
			rt = expect
		}
		return &Tuple{
			Type:  rt,
			Elems: elems,
		}, nil
	case ast.Array != nil:
		if len(ast.Array.Exprs) == 0 {
			if expect == nil || !IsArrayTypeAndSon(expect) {
				return nil, utils.Errorf(ast.Position, "expect a array type")
			}
			return &EmptyArray{Type: expect}, nil
		}
		expects := make([]Type, len(ast.Array.Exprs))
		if expect != nil {
			if at, ok := GetBaseType(expect).(*TypeArray); ok && at.Size == uint(len(ast.Array.Exprs)) {
				for i := range expects {
					expects[i] = at.Elem
				}
			}
		}
		elems, err := analyseExprList(ctx, expects, *ast.Array)
		if err != nil {
			return nil, err
		}
		var errors []utils.Error
		for i, e := range ast.Array.Exprs {
			elems[i], err = expectExpr(e.Position, elems[0].GetType(), elems[i])
			if err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) == 1 {
			return nil, errors[0]
		} else if len(errors) > 1 {
			return nil, utils.NewMultiError(errors...)
		}
		var rt Type = NewArrayType(uint(len(elems)), elems[0].GetType())
		if expect != nil && GetDepthBaseType(expect).Equal(GetDepthBaseType(rt)) {
			rt = expect
		}
		return &Array{
			Type:  rt,
			Elems: elems,
		}, nil
	case ast.Struct != nil:
		if len(ast.Struct.Exprs) == 0 {
			if expect == nil || !IsStructTypeAndSon(expect) {
				return nil, utils.Errorf(ast.Position, "expect a struct type")
			}
			return &EmptyStruct{Type: expect}, nil
		}
		if expect == nil || !IsStructTypeAndSon(expect) {
			return nil, utils.Errorf(ast.Position, "expect a struct type")
		} else if GetBaseType(expect).(*TypeStruct).Fields.Length() != len(ast.Struct.Exprs) {
			return nil, utils.Errorf(ast.Position, "expect `%d` fields", len(ast.Struct.Exprs))
		}
		expects := make([]Type, len(ast.Struct.Exprs))
		for iter := GetBaseType(expect).(*TypeStruct).Fields.Begin(); iter.HasValue(); iter.Next() {
			expects[iter.Index()] = iter.Value()
		}
		fields, err := analyseExprList(ctx, expects, *ast.Struct)
		if err != nil {
			return nil, err
		}
		for i, e := range fields {
			expects[i] = e.GetType()
		}
		return &Struct{
			Type:   expect,
			Fields: fields,
		}, nil
	default:
		panic("")
	}
}

// 常量表达式
func analyseConstant(expect Type, ast parse.Constant) (Expr, utils.Error) {
	switch {
	case ast.Null != nil:
		if expect == nil || (!IsPtrTypeAndSon(expect) && !IsFuncTypeAndSon(expect)) {
			return nil, utils.Errorf(ast.Position, "expect a pointer type")
		}
		return &Null{Type: expect}, nil
	case ast.Int != nil:
		if expect == nil || !IsNumberTypeAndSon(expect) {
			expect = Isize
		}
		if IsIntTypeAndSon(expect) {
			return &Integer{
				Type:  expect,
				Value: *ast.Int,
			}, nil
		} else {
			return &Float{
				Type:  expect,
				Value: float64(*ast.Int),
			}, nil
		}
	case ast.Float != nil:
		if expect == nil || !IsFloatTypeAndSon(expect) {
			expect = F64
		}
		return &Float{
			Type:  expect,
			Value: *ast.Float,
		}, nil
	case ast.Bool != nil:
		if expect == nil || !IsBoolTypeAndSon(expect) {
			expect = Bool
		}
		return &Boolean{
			Type:  expect,
			Value: bool(*ast.Bool),
		}, nil
	case ast.Char != nil:
		return &Integer{
			Type:  I32,
			Value: int64(*ast.Char),
		}, nil
	case ast.CString != nil:
		elems := make([]Expr, len(*ast.CString))
		for i, e := range *ast.CString {
			elems[i] = &Integer{
				Type:  I8,
				Value: int64(e),
			}
		}
		return &Array{
			Type:  NewArrayType(uint(len(elems)), I8),
			Elems: elems,
		}, nil
	case ast.String != nil:
		elems := make([]Expr, len(*ast.String))
		for i, e := range *ast.String {
			elems[i] = &Integer{
				Type:  I32,
				Value: int64(e),
			}
		}
		return &Array{
			Type:  NewArrayType(uint(len(elems)), I32),
			Elems: elems,
		}, nil
	default:
		panic("")
	}
}

// 常量表达式
func analyseConstantExpr(expect Type, ast parse.Expr) (Expr, utils.Error) {
	if ast.Assign.Suffix == nil &&
		len(ast.Assign.Left.Next) == 0 &&
		len(ast.Assign.Left.Left.Next) == 0 &&
		len(ast.Assign.Left.Left.Left.Next) == 0 &&
		len(ast.Assign.Left.Left.Left.Left.Next) == 0 &&
		len(ast.Assign.Left.Left.Left.Left.Left.Next) == 0 &&
		len(ast.Assign.Left.Left.Left.Left.Left.Left.Suffix) == 0 &&
		ast.Assign.Left.Left.Left.Left.Left.Left.Unary.Opera == nil &&
		len(ast.Assign.Left.Left.Left.Left.Left.Left.Unary.Postfix.Suffix) == 0 {
		primary := ast.Assign.Left.Left.Left.Left.Left.Left.Unary.Postfix.Primary
		switch {
		case primary.Constant != nil:
			return analyseConstant(expect, *primary.Constant)
		case primary.Array != nil:
			if len(primary.Array.Exprs) == 0 {
				if expect == nil || !IsArrayTypeAndSon(expect) {
					return nil, utils.Errorf(primary.Position, "expect a array type")
				}
				return &EmptyArray{Type: expect}, nil
			}
			expects := make([]Type, len(primary.Array.Exprs))
			if expect != nil {
				if at, ok := GetBaseType(expect).(*TypeArray); ok && at.Size == uint(len(primary.Array.Exprs)) {
					for i := range expects {
						expects[i] = at.Elem
					}
				}
			}
			elems, err := analyseConstantExprList(expects, *primary.Array)
			if err != nil {
				return nil, err
			}
			var errors []utils.Error
			for i, e := range primary.Array.Exprs {
				elems[i], err = expectExpr(e.Position, elems[0].GetType(), elems[i])
				if err != nil {
					errors = append(errors, err)
				}
			}
			if len(errors) == 1 {
				return nil, errors[0]
			} else if len(errors) > 1 {
				return nil, utils.NewMultiError(errors...)
			}
			var rt Type = NewArrayType(uint(len(elems)), elems[0].GetType())
			if expect != nil && GetDepthBaseType(expect).Equal(GetDepthBaseType(rt)) {
				rt = expect
			}
			return &Array{
				Type:  rt,
				Elems: elems,
			}, nil
		case primary.Tuple != nil:
			if len(primary.Tuple.Exprs) == 0 {
				if expect == nil || !IsTupleTypeAndSon(expect) {
					return nil, utils.Errorf(primary.Position, "expect a tuple type")
				}
				return &EmptyTuple{Type: expect}, nil
			} else if len(primary.Tuple.Exprs) == 1 && (expect == nil || !IsTupleTypeAndSon(expect) || len(GetBaseType(expect).(*TypeTuple).Elems) != 1) {
				return analyseConstantExpr(expect, primary.Tuple.Exprs[0])
			}
			expects := make([]Type, len(primary.Tuple.Exprs))
			if expect != nil {
				if tt, ok := GetBaseType(expect).(*TypeTuple); ok && len(tt.Elems) == len(primary.Tuple.Exprs) {
					for i := range expects {
						expects[i] = tt.Elems[i]
					}
				}
			}
			elems, err := analyseConstantExprList(expects, *primary.Tuple)
			if err != nil {
				return nil, err
			}
			for i, e := range elems {
				expects[i] = e.GetType()
			}
			var rt Type = NewTupleType(expects...)
			if expect != nil && GetDepthBaseType(expect).Equal(GetDepthBaseType(rt)) {
				rt = expect
			}
			return &Tuple{
				Type:  rt,
				Elems: elems,
			}, nil
		case primary.Struct != nil:
			if len(primary.Struct.Exprs) == 0 {
				if expect == nil || !IsStructTypeAndSon(expect) {
					return nil, utils.Errorf(primary.Position, "expect a struct type")
				}
				return &EmptyStruct{Type: expect}, nil
			}
			if expect == nil || !IsStructTypeAndSon(expect) {
				return nil, utils.Errorf(primary.Position, "expect a struct type")
			} else if GetBaseType(expect).(*TypeStruct).Fields.Length() != len(primary.Struct.Exprs) {
				return nil, utils.Errorf(ast.Position, "expect `%d` fields", len(primary.Struct.Exprs))
			}
			expects := make([]Type, len(primary.Struct.Exprs))
			for iter := GetBaseType(expect).(*TypeStruct).Fields.Begin(); iter.HasValue(); iter.Next() {
				expects[iter.Index()] = iter.Value()
			}
			fields, err := analyseConstantExprList(expects, *primary.Struct)
			if err != nil {
				return nil, err
			}
			for i, e := range fields {
				expects[i] = e.GetType()
			}
			return &Struct{
				Type:   expect,
				Fields: fields,
			}, nil
		}
	}
	return nil, utils.Errorf(ast.Position, "expect a constant value")
}

// 期待指定类型的表达式
func expectExpr(pos utils.Position, expect Type, expr Expr) (Expr, utils.Error) {
	if !expr.GetType().Equal(expect) {
		return nil, utils.Errorf(pos, "expect type `%s`", expect)
	}
	return expr, nil
}

// 期待指定类型的表达式及其子类型
func expectExprAndSon(pos utils.Position, expect Type, expr Expr) (Expr, utils.Error) {
	if !GetDepthBaseType(expr.GetType()).Equal(GetDepthBaseType(expect)) {
		return nil, utils.Errorf(pos, "expect type `%s`", expect)
	}
	return expr, nil
}

// 获取类型默认值
func getDefaultExprByType(t Type) Expr {
	switch GetBaseType(t).(type) {
	case *typeBasic:
		switch {
		case IsNoneType(t):
			panic("")
		case IsIntType(t):
			return &Integer{
				Type:  t,
				Value: 0,
			}
		case IsFloatType(t):
			return &Float{
				Type:  t,
				Value: 0,
			}
		case IsBoolType(t):
			return &Boolean{
				Type:  t,
				Value: false,
			}
		default:
			panic("")
		}
	case *TypeFunc:
		return &Null{Type: t}
	case *TypeArray:
		return &EmptyArray{Type: t}
	case *TypeTuple:
		return &EmptyTuple{Type: t}
	case *TypeStruct:
		return &EmptyStruct{Type: t}
	case *TypePtr:
		return &Null{Type: t}
	default:
		panic("")
	}
}

// 表达式列表
func analyseExprList(ctx *blockContext, expects []Type, ast parse.ExprList) ([]Expr, utils.Error) {
	exprs := make([]Expr, len(ast.Exprs))
	var errors []utils.Error
	for i, e := range ast.Exprs {
		var expect Type
		if len(expects) == len(ast.Exprs) {
			expect = expects[i]
		}
		expr, err := analyseExpr(ctx, expect, e)
		if err != nil {
			errors = append(errors, err)
		} else {
			exprs[i] = expr
		}
	}

	if len(errors) == 0 {
		return exprs, nil
	} else if len(errors) == 1 {
		return nil, errors[0]
	} else {
		return nil, utils.NewMultiError(errors...)
	}
}

// 常量表达式列表
func analyseConstantExprList(expects []Type, ast parse.ExprList) ([]Expr, utils.Error) {
	exprs := make([]Expr, len(ast.Exprs))
	var errors []utils.Error
	for i, e := range ast.Exprs {
		var expect Type
		if len(expects) == len(ast.Exprs) {
			expect = expects[i]
		}
		expr, err := analyseConstantExpr(expect, e)
		if err != nil {
			errors = append(errors, err)
		} else {
			exprs[i] = expr
		}
	}

	if len(errors) == 0 {
		return exprs, nil
	} else if len(errors) == 1 {
		return nil, errors[0]
	} else {
		return nil, utils.NewMultiError(errors...)
	}
}

// 标识符
func analyseIdent(ctx *blockContext, ast parse.Ident) (Expr, utils.Error) {
	if ast.Package == nil {
		v := ctx.GetValue(ast.Name.Value)
		if v == nil {
			return nil, utils.Errorf(ast.Position, "unknown identifier")
		}
		return v, nil
	} else {
		pkg := ctx.GetPackageContext().externs[ast.Package.Value]
		if pkg == nil {
			return nil, utils.Errorf(ast.Package.Position, "unknown `%s`", ast.Package.Value)
		}
		value := pkg.GetValue(ast.Name.Value)
		if !value.First || value.Second == nil {
			return nil, utils.Errorf(ast.Name.Position, "unknown `%s`", ast.Name.Value)
		}
		return value.Second, nil
	}
}
