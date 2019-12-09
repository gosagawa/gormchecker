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
		(*ast.CallExpr)(nil),
	}
	functions := make(map[int][]string)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		position := pass.Fset.Position(n.Pos())
		fmt.Printf("### pos %v\n", n.Pos())
		fmt.Println(position)
		ast.Print(pass.Fset, n)
		switch n := n.(type) {
		case *ast.CallExpr:
			switch f := n.Fun.(type) {
			case *ast.SelectorExpr:
				functionName := f.Sel.Name
				functions[position.Line] = append(functions[position.Line], functionName)
			}
		}

	})

	fmt.Println(functions)
	return nil, nil
}
