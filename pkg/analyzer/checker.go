package analyzer

import (
	"go/ast"
	"go/types"
	"selectellinter/config"
	"selectellinter/pkg/analyzer/detector"

	"golang.org/x/tools/go/analysis"
)

var knownPacks = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}

func runWithConfig(cfg *config.Config) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		detectors := []*detector.AWSDetector{}
		if cfg.Rules.CustomPatterns {
			for _, customPattern := range cfg.CustomPatterns {
				pattern := customPattern.Pattern
				detector := detector.NewAWSDetector(pattern)
				detectors = append(detectors, detector)
			}
		}
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkgIdent, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				obj := pass.TypesInfo.ObjectOf(pkgIdent)
				if obj == nil {
					return true
				}
				var isLog bool
				switch obj := obj.(type) {
				case *types.PkgName:
					isLog = isLogger(obj)
				case *types.Var:
					isLog = isLoggerVar(obj.Type())
				case *types.Func:
					isLog = isLoggerFunc(obj)
				}

				if isLog {
					checkLogMessage(pass, call, cfg, detectors)
				}

				return true
			})
		}
		return nil, nil
	}
}

// Check if the call expretion is a log call
func isLogger(pkgName *types.PkgName) bool {
	pkgPath := pkgName.Imported().Path()

	if _, ok := knownPacks[pkgPath]; ok {
		return true
	}
	return false
}

// Check if the call expretion is a logger variable
func isLoggerVar(t types.Type) bool {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	pkg := named.Obj().Pkg()
	if pkg == nil {
		return false
	}
	//typeName := named.Obj().Name()
	pkgPath := pkg.Path()
	if _, ok := knownPacks[pkgPath]; ok {
		return true
	}
	return false
}

// Check if the call expretion is a function that returns a logger variable
func isLoggerFunc(fn *types.Func) bool {
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}

	res := sig.Results()
	for i := range res.Len() {
		if isLoggerVar(res.At(i).Type()) {
			return true
		}
	}
	return false
}
