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
	f, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		log.Fatalf("ParseFile failed: %s", err)
	}

	fmt.Printf("%#v", f)
}
