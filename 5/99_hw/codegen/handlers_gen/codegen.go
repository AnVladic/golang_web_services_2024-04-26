package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "api.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create("api_handlers.go")
	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "net/url"`)
	//fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out) // empty line

	structGenerator(out, node)

	methodsMap := funcRequires(node)
	funcGenerator(out, methodsMap)
}
