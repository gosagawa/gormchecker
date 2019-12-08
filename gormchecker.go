package gormchecker

import (
	"go/ast"
	"go/token"

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
		fset := token.NewFileSet()
		ast.Print(fset, n)

		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "Gopher" {
				pass.Reportf(n.Pos(), "name of identifier must not be ’Gopher’")
			}
		}
	})

	return nil, nil
}
