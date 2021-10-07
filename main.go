package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
)

type visitor struct {
	fset   *token.FileSet
	prog   *loader.PackageInfo
	fmtpkg *types.Package
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	if c, ok := node.(*ast.CallExpr); ok {
		if s, ok := c.Fun.(*ast.SelectorExpr); ok {
			if v.prog.Info.ObjectOf(s.Sel).Pkg() == v.fmtpkg {
				if s.Sel.Name == "Printf" || s.Sel.Name == "Println" {
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

	main := prog.Package(astf.Name.Name)
	fmtpkg := prog.Package("fmt").Pkg
	v := &visitor{l.Fset, main, fmtpkg}
	ast.Walk(v, astf)
}
