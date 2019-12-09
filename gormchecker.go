package gormchecker

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is analizer for gormchecker
var Analyzer = &analysis.Analyzer{
	Name: "gormchecker",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// Doc is explain for gormchecker
const Doc = "gormchecker is analyzer for gorm using funtion"

func run(pass *analysis.Pass) (interface{}, error) {

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fmt.Printf("### pos %v\n", n.Pos())
		fmt.Println(pass.Fset.Position(n.Pos()))
		ast.Print(pass.Fset, n)

		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "Gopher" {
				pass.Reportf(n.Pos(), "name of identifier must not be ’Gopher’")
			}
		}
	})

	return nil, nil
}
