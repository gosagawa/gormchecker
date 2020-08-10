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
		(*ast.CallExpr)(nil),
		(*ast.FuncDecl)(nil),
	}
	functions := make(map[string]map[int][]string)
	includeFunctions := make(map[string]map[string]int)
	var baseFunction string
	baseFunctionPos := make(map[string]token.Pos)

	inspect.Preorder(nodeFilter, func(n ast.Node) {

		position := pass.Fset.Position(n.Pos())
		switch n := n.(type) {
		case *ast.FuncDecl:
			baseFunction = n.Name.Name
			baseFunctionPos[baseFunction] = n.Pos()
		case *ast.CallExpr:
			switch f := n.Fun.(type) {
			case *ast.SelectorExpr:
				if x, ok := f.X.(*ast.Ident); ok && x.Name != "db" {
					break
				}
				functionName := f.Sel.Name
				if _, ok := functions[position.Filename]; !ok {
					functions[position.Filename] = make(map[int][]string)
				}

				if x, ok := f.X.(*ast.Ident); ok && x.Name == "db" {
					if _, ok := includeFunctions[baseFunction]; !ok {
						includeFunctions[baseFunction] = make(map[string]int)
					}
					if _, ok := includeFunctions[baseFunction][functionName]; !ok {
						includeFunctions[baseFunction][functionName] = 0
					}
					includeFunctions[baseFunction][functionName]++
				}

				functions[position.Filename][position.Line] = append(functions[position.Filename][position.Line], functionName)
				if len(functions[position.Filename][position.Line]) > 1 {
					pass.Reportf(n.Pos(), "do not use pipe")
				}
			}
		}
	})

	for baseFunction, v := range includeFunctions {
		findCount := 0
		firstCount := 0
		if i, ok := v["Find"]; ok {
			findCount = i
		}
		if i, ok := v["First"]; ok {
			firstCount = i
		}
		if findCount == 0 && firstCount == 0 {
			pass.Reportf(baseFunctionPos[baseFunction], "not have Find or First")
		}
		if findCount+firstCount > 1 {
			pass.Reportf(baseFunctionPos[baseFunction], "have two more Find or First")
		}
	}

	return nil, nil
}
