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

		file := pass.Fset.File(n.Pos())
		position := pass.Fset.Position(n.Pos())
		switch n := n.(type) {
		case *ast.FuncDecl:
			baseFunction = file.Name() + "_" + n.Name.Name
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
		createCount := 0
		updateCount := 0
		deleteCount := 0
		whereCount := 0
		queryFunctionCount := 0
		queryExecuteFunctionCount := 0
		selectExecuteFunctionCount := 0
		isEditFunction := false
		if i, ok := v["Find"]; ok {
			queryFunctionCount++
			queryExecuteFunctionCount += i
			selectExecuteFunctionCount += i
		}
		if i, ok := v["First"]; ok {
			queryFunctionCount += i
			queryExecuteFunctionCount += i
			selectExecuteFunctionCount += i
		}
		if i, ok := v["Pluck"]; ok {
			queryFunctionCount += i
			queryExecuteFunctionCount += i
			selectExecuteFunctionCount += i
		}
		if i, ok := v["Scan"]; ok {
			queryFunctionCount += i
			queryExecuteFunctionCount += i
			selectExecuteFunctionCount += i
		}
		if i, ok := v["Count"]; ok {
			queryFunctionCount += i
			queryExecuteFunctionCount += i
			selectExecuteFunctionCount += i
		}
		if i, ok := v["Create"]; ok {
			createCount = i
			queryFunctionCount += i
			queryExecuteFunctionCount += i
		}
		if i, ok := v["Update"]; ok {
			updateCount = i
			queryFunctionCount += i
			isEditFunction = true
			queryExecuteFunctionCount += i
		}
		if i, ok := v["Updates"]; ok {
			updateCount = i
			queryFunctionCount += i
			isEditFunction = true
			queryExecuteFunctionCount += i
		}
		if i, ok := v["Delete"]; ok {
			deleteCount = i
			queryFunctionCount++
			isEditFunction = true
			queryExecuteFunctionCount++
		}
		if i, ok := v["Where"]; ok {
			whereCount = i
			queryFunctionCount++
		}
		if queryFunctionCount == 0 {
			continue
		}
		if queryExecuteFunctionCount == 0 {
			pass.Reportf(baseFunctionPos[baseFunction], "not have query execution function like Find, Create, Update")
		}
		if selectExecuteFunctionCount > 1 {
			pass.Reportf(baseFunctionPos[baseFunction], "have two more select function like Find, First")
		}
		if createCount > 1 {
			pass.Reportf(baseFunctionPos[baseFunction], "have two more Create")
		}
		if updateCount > 1 {
			pass.Reportf(baseFunctionPos[baseFunction], "have two more Update")
		}
		if deleteCount > 1 {
			pass.Reportf(baseFunctionPos[baseFunction], "have two more Delete")
		}
		if isEditFunction && whereCount == 0 {
			pass.Reportf(baseFunctionPos[baseFunction], "no where Edit function")
		}
	}

	return nil, nil
}
