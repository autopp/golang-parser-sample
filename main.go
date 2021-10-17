package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

type visitor struct {
	fset     *token.FileSet
	typeInfo *types.Info
	fmtpkg   *types.Package
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	if c, ok := node.(*ast.CallExpr); ok {
		if s, ok := c.Fun.(*ast.SelectorExpr); ok {
			if v.typeInfo.ObjectOf(s.Sel).Pkg() == v.fmtpkg {
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
	cfg := &packages.Config{
		Tests: false,
		Mode:  packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, os.Args[1])
	if err != nil {
		log.Fatalf("packages.Load failed: %s", err)
	}

	pkg := pkgs[0]
	var fmtPkg *types.Package
	for _, p := range pkg.Types.Imports() {
		if p.Path() == "fmt" {
			fmtPkg = p
			break
		}
	}

	if fmtPkg == nil {
		return
	}

	v := &visitor{fset: pkg.Fset, typeInfo: pkg.TypesInfo, fmtpkg: fmtPkg}
	ast.Walk(v, pkg.Syntax[0])
}
