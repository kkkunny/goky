package codegen

import (
	"fmt"
	"github.com/kkkunny/klang/src/compiler/internal/analyse"
	"strings"
)

// 类型
func (self *CodeGenerator) generateType(mean analyse.Type) string {
	switch typ := mean.(type) {
	case *analyse.TypeFunc:
		key := typ.String()
		if name, ok := self.types[key]; ok {
			return name
		}
		ret := self.generateType(typ.Ret)
		params := make([]string, len(typ.Params))
		for i, p := range typ.Params {
			params[i] = self.generateType(p)
		}
		name := fmt.Sprintf("__t%d", self.typeCount)
		self.typeCount++
		self.typedef.WriteString(fmt.Sprintf("typedef %s(*%s)(%s);\n", ret, name, strings.Join(params, "")))
		self.types[key] = name
		return name
	case *analyse.TypeArray:
		key := typ.String()
		if name, ok := self.types[key]; ok {
			return name
		}
		elem := self.generateType(typ.Elem)
		name := fmt.Sprintf("__t%d", self.typeCount)
		self.typeCount++
		self.typedef.WriteString(fmt.Sprintf("typedef struct{%s data[%d];} %s;\n", elem, typ.Size, name))
		self.types[key] = name
		return name
	case *analyse.TypeTuple:
		key := typ.String()
		if name, ok := self.types[key]; ok {
			return name
		}
		elems := make([]string, len(typ.Elems))
		for i, e := range typ.Elems {
			elems[i] = fmt.Sprintf("%s e%d;", self.generateType(e), i)
		}
		name := fmt.Sprintf("__t%d", self.typeCount)
		self.typeCount++
		self.typedef.WriteString(fmt.Sprintf("typedef struct{%s} %s;\n", strings.Join(elems, ""), name))
		self.types[key] = name
		return name
	case *analyse.TypeStruct:
		key := typ.String()
		if name, ok := self.types[key]; ok {
			return name
		}
		elems := make([]string, typ.Fields.Length())
		for iter := typ.Fields.Begin(); iter.HasValue(); iter.Next() {
			elems[iter.Index()] = fmt.Sprintf("%s f%d;", self.generateType(iter.Value()), iter.Index())
		}
		name := fmt.Sprintf("__t%d", self.typeCount)
		self.typeCount++
		self.typedef.WriteString(fmt.Sprintf("typedef struct{%s} %s;\n", strings.Join(elems, ""), name))
		self.types[key] = name
		return name
	case *analyse.TypePtr:
		return self.generateType(typ.Elem) + "*"
	default:
		switch {
		case analyse.IsNoneType(typ):
			return "void"
		case analyse.IsSintType(typ):
			switch typ {
			case analyse.I8:
				return "char"
			case analyse.I16:
				return "short"
			case analyse.I32:
				return "int"
			case analyse.I64:
				return "long"
			case analyse.Isize:
				return "long"
			default:
				panic("")
			}
		case analyse.IsUintType(typ):
			switch typ {
			case analyse.U8:
				return "unsigned char"
			case analyse.U16:
				return "unsigned short"
			case analyse.U32:
				return "unsigned int"
			case analyse.U64:
				return "unsigned long"
			case analyse.Usize:
				return "unsigned long"
			default:
				panic("")
			}
		case analyse.IsFloatType(typ):
			switch typ {
			case analyse.F32:
				return "float"
			case analyse.F64:
				return "double"
			default:
				panic("")
			}
		case analyse.IsBoolType(typ):
			return "char"
		default:
			panic("")
		}
	}
}
