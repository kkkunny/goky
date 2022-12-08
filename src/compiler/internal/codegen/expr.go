package codegen

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/analyse"
	stlutil "github.com/kkkunny/stl/util"
	"strings"
)

// 表达式
func (self *CodeGenerator) generateExpr(mean analyse.Expr) string {
	switch expr := mean.(type) {
	case *analyse.Null, *analyse.Integer, *analyse.Float, *analyse.EmptyStruct, *analyse.EmptyArray, *analyse.EmptyTuple:
		return self.generateConstantExpr(mean)
	case *analyse.Binary:
		l, r := self.generateExpr(expr.Left), self.generateExpr(expr.Right)
		switch expr.Opera {
		case "+", "-", "*", "/", "&", "|", "^", "<<", ">>", "&&", "||":
			return fmt.Sprintf("(%s) %s (%s)", l, expr.Opera, r)
		case "%":
			self.include["math"] = struct{}{}
			return fmt.Sprintf("fmod(%s, %s)", l, r)
		default:
			panic("")
		}
	case *analyse.Variable:
		return self.vars[expr]
	case *analyse.Function:
		return self.vars[expr]
	case *analyse.Call:
		fun := self.generateExpr(expr.Func)
		args := make([]string, len(expr.Args))
		for i, a := range expr.Args {
			args[i] = self.generateExpr(a)
		}
		return fmt.Sprintf("%s(%s)", fun, strings.Join(args, ", "))
	case *analyse.Param:
		return self.vars[expr]
	case *analyse.Array:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	case *analyse.Assign:
		return fmt.Sprintf("%s %s %s", self.generateExpr(expr.Left), expr.Opera, self.generateExpr(expr.Right))
	case *analyse.Boolean:
		return stlutil.Ternary(expr.Value, "true", "false")
	case *analyse.Equal:
		l, r := self.generateExpr(expr.Left), self.generateExpr(expr.Right)
		return fmt.Sprintf("(%s) %s (%s)", l, expr.Opera, r)
	case *analyse.Unary:
		return fmt.Sprintf("%s(%s)", expr.Opera, self.generateExpr(expr.Value))
	case *analyse.Index:
		fromType := expr.From.GetType()
		f, i := self.generateExpr(expr.From), self.generateExpr(expr.Index)
		switch {
		case analyse.IsArrayType(fromType):
			return fmt.Sprintf("(%s).data[%s]", f, i)
		case analyse.IsPtrType(fromType):
			return fmt.Sprintf("(%s)[%s]", f, i)
		case analyse.IsTupleType(fromType):
			return fmt.Sprintf("(%s).e%s", f, i)
		default:
			panic("")
		}
	case *analyse.Select:
		// TODO
		panic("")
	case *analyse.Tuple:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	case *analyse.Struct:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Fields))
		for i, e := range expr.Fields {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	case *analyse.GetField:
		var index uint
		for iter := expr.From.GetType().(*analyse.TypeStruct).Fields.Begin(); iter.HasValue(); iter.Next() {
			if iter.Key() == expr.Index {
				index = uint(iter.Index())
				break
			}
		}
		return fmt.Sprintf("(%s).f%d", self.generateExpr(expr.From), index)
	case *analyse.Covert:
		return fmt.Sprintf("(%s)(%s)", self.generateType(expr.To), self.generateExpr(expr.From))
	case *analyse.GlobalVariable:
		return self.vars[expr]
	default:
		panic("")
	}
}

// 常量表达式
func (self *CodeGenerator) generateConstantExpr(mean analyse.Expr) string {
	switch expr := mean.(type) {
	case *analyse.Null:
		return "NULL"
	case *analyse.Integer:
		return fmt.Sprintf("%d", expr.Value)
	case *analyse.Float:
		return fmt.Sprintf("%f", expr.Value)
	case *analyse.EmptyArray, *analyse.EmptyTuple, *analyse.EmptyStruct:
		return "{}"
	case *analyse.Array:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	case *analyse.Tuple:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Elems))
		for i, e := range expr.Elems {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	case *analyse.Struct:
		t := self.generateType(expr.GetType())
		elems := make([]string, len(expr.Fields))
		for i, e := range expr.Fields {
			elems[i] = self.generateExpr(e)
		}
		return fmt.Sprintf("(%s){%s}", t, strings.Join(elems, ", "))
	default:
		panic("")
	}
}

// 比较
func (self *CodeGenerator) equal(ast analyse.Type, l string, r string) string {
	switch {
	case analyse.IsArrayType(ast):
		// TODO
		panic("")
	case analyse.IsTupleType(ast):
		// TODO
		panic("")
	case analyse.IsStructType(ast):
		// TODO
		panic("")
	default:
		return fmt.Sprintf("(%s) == (%s)", l, r)
	}
}
