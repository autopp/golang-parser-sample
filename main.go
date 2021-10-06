package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
)

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	if c, ok := node.(*ast.CallExpr); ok {
		if s, ok := c.Fun.(*ast.SelectorExpr); ok {
			if p, ok := s.X.(*ast.Ident); ok {
				if p.Name == "fmt" && (s.Sel.Name == "Printf" || s.Sel.Name == "Println") {
					buf := &bytes.Buffer{}
					printer.Fprint(buf, v.fset, node)
					fmt.Printf("%s: %s\n", v.fset.Position(node.Pos()), buf.String())
				}
			}
		}
	}
	return v
}

func main() {
	l := loader.Config{ParserMode: parser.ParseComments}
	astf, err := l.ParseFile(os.Args[1], nil)
	if err != nil {
		log.Fatalf("ParseFile failed: %s", err)
	}
	l.CreateFromFiles("", astf)
	prog, err := l.Load()
	if err != nil {
		log.Fatalf("Load failed: %s", err)
	}

	v := &visitor{l.Fset}
	main := prog.Package("main")
	for _, f := range main.Files {
		ast.Walk(v, f)
	}
}
