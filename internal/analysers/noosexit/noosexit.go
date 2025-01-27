package noosexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Analyzer — определение анализатора
var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "запрещает использовать os.Exit в функции main",
	Run:  run,
}

// run — основная функция анализа
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем вызовы os.Exit
			if call, ok := n.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					if fun.Sel.Name == "Exit" {
						if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "os" {
							pass.Reportf(call.Pos(), "использование os.Exit запрещено в функции main")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
