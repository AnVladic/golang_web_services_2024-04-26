package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"regexp"
)

func funcRequires(node *ast.File) *map[string][]*ast.FuncDecl {
	var methodsMap = make(map[string][]*ast.FuncDecl)

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			fmt.Printf("SKIP %#T is not *ast.FuncDecl\n", funcDecl)
			continue
		}

		if funcDecl.Doc == nil {
			fmt.Printf("SKIP %#T is not *ast.FuncDecl.Doc\n", funcDecl)
			continue
		}

		if funcDecl.Recv == nil || len(funcDecl.Recv.List) != 1 {
			fmt.Printf("SKIP %#T is not *ast.FuncDecl.Recv\n", funcDecl)
			continue
		}

		firstReceiver := funcDecl.Recv.List[0]
		receiverType, ok := firstReceiver.Type.(*ast.StarExpr)
		if !ok {
			fmt.Printf("SKIP %#T is not *ast.StarExpr\n", funcDecl)
			continue
		}

		ident, ok := receiverType.X.(*ast.Ident)
		if !ok {
			fmt.Printf("SKIP %#T is not *ast.Ident\n", receiverType)
			continue
		}

		methodsList, ok := methodsMap[ident.Name]
		if !ok {
			methodsList = []*ast.FuncDecl{funcDecl}
		} else {
			methodsList = append(methodsList, funcDecl)
		}
		methodsMap[ident.Name] = methodsList
	}
	return &methodsMap
}

func funcGenerator(out *os.File, methodsMap *map[string][]*ast.FuncDecl) {
	for structName, methods := range *methodsMap {
		fmt.Fprintf(out, `func (h *%s) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
`, structName)
		for _, method := range methods {
			re := regexp.MustCompile(`apigen:api\s+({.*})`)
			match := re.FindStringSubmatch(method.Doc.Text())

			if len(match) < 0 {
				continue
			}
			var result map[string]interface{}
			err := json.Unmarshal([]byte(match[1]), &result)
			if err != nil {
				continue
			}
			requestMethod := result["method"]
			if requestMethod == nil {
				requestMethod = ""
			}
			_, _ = fmt.Fprintf(out, `	case "%s":
		requestValues, err := validRequest(w, r, "%v", %v)
		if err != nil {
			MarshalAndWrite(w, err)
			return
		}
`, result["url"], requestMethod, result["auth"])
			for _, param := range method.Type.Params.List[1:] {
				fmt.Fprintf(out, `		param := %v{}
		if err := param.Valid(requestValues); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			MarshalAndWrite(w, err)
			return
		}
		response, err := h.%s(r.Context(), param)
		if err != nil {
			SetFuncError(w, err)
			return
		}
		MarshalAndWrite(w, &ResponseError{"", response})
`, param.Type, method.Name)
			}
		}
		_, _ = fmt.Fprint(out, `
	default:
		w.WriteHeader(http.StatusNotFound)
		MarshalAndWrite(w, &ResponseError{"unknown method", nil})
	}
}

`)
	}
}
