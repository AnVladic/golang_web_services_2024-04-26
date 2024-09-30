package main

import (
	"fmt"
	"go/ast"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func structGenerator(out *os.File, node *ast.File) {
	for _, decl := range node.Decls {
		g, ok := decl.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %#T is not *ast.GenDecl\n", decl)
			continue
		}
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				fmt.Printf("SKIP %#T is not ast.StructType\n", currStruct)
				continue
			}

			isRequire := false
			for _, field := range currStruct.Fields.List {
				if field.Tag != nil {
					tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
					apivalidator := tag.Get("apivalidator")
					if apivalidator == "-" || apivalidator == "" {
						continue
					}
				} else {
					continue
				}
				isRequire = true
				break
			}
			if isRequire {
				generateValidMethod(out, currStruct, currType)
			}
		}
	}
}

func convertEnumsToString(enums []string) string {
	if len(enums) == 0 {
		return "nil"
	}

	quotedEnums := make([]string, len(enums))
	for i, v := range enums {
		quotedEnums[i] = fmt.Sprintf("\"%s\"", v)
	}

	return "[]string{" + strings.Join(quotedEnums, ", ") + "}"
}

func convertEnumsToIntString(enums []string) string {
	if len(enums) == 0 {
		return "nil"
	}

	indexes := make([]string, len(enums))
	for i := range enums {
		indexes[i] = strconv.Itoa(i)
	}

	return "[]int{" + strings.Join(indexes, ", ") + "}"
}

func generateValidMethod(out *os.File, currStruct *ast.StructType, currType *ast.TypeSpec) {
	fmt.Fprintf(out, `func (s *%s) Valid(query url.Values) error {
	var err error
`, currType.Name.Name)
	for _, field := range currStruct.Fields.List {
		fieldType, ok := field.Type.(*ast.Ident)
		if !ok {
			fmt.Printf("SKIP %s %#T is not *ast.Ident\n", field.Type, field.Type)
			continue
		}

		name := field.Names[0].Name
		validator := field.Tag.Value[15 : len(field.Tag.Value)-2]
		validators := strings.Split(validator, ",")
		isRequire := false
		paramName := strings.ToLower(name)
		minValue := math.MinInt
		maxValue := math.MaxInt
		defualt := ""
		var enums []string
		for _, rule := range validators {
			switch {
			case strings.HasPrefix(rule, "required"):
				isRequire = true
			case strings.Contains(rule, "paramname="):
				paramName = strings.Replace(rule, "paramname=", "", 1)
			case strings.Contains(rule, "enum="):
				for _, enum := range strings.Split(rule[5:], "|") {
					enums = append(enums, strings.TrimSpace(enum))
				}
			case strings.Contains(rule, "min="):
				minValue, _ = strconv.Atoi(rule[4:])
			case strings.Contains(rule, "max="):
				maxValue, _ = strconv.Atoi(rule[4:])
			case strings.Contains(rule, "default="):
				defualt = rule[8:]
			}
		}
		//validators := strings.Split(validator, ",")

		switch fieldType.Name {
		case "int":
			fmt.Fprintf(out, "	if s.%v, err = ValidInt("+
				"\n\t\tquery, "+
				"\n\t\t\"%s\", "+
				"\n\t\t%v, "+
				"\n\t\t%v,"+
				"\n\t\t%d,"+
				"\n\t\t%d,"+
				"\n\t\t); err != nil {\n",
				name, paramName, isRequire, convertEnumsToIntString(enums), minValue, maxValue)
		default:
			fmt.Fprintf(out, "	if s.%v, err = ValidString("+
				"\n\t\tquery, "+
				"\n\t\t\"%s\", "+
				"\n\t\t%v, "+
				"\n\t\t%v,"+
				"\n\t\t%d,"+
				"\n\t\t%d,"+
				"\n\t\t\"%s\","+
				"\n\t); err != nil {\n",
				name, paramName, isRequire, convertEnumsToString(enums), minValue, maxValue, defualt)
		}
		fmt.Fprint(out, "\t\treturn err\n\t}\n")
	}
	fmt.Fprint(out, "	return nil\n}\n\n")
}
