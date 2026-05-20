package analyzer

import (
	"go/ast"
	"go/types"
	"linter/config"
	"linter/pkg/analyzer/detector"

	"golang.org/x/tools/go/analysis"
)

var LogMethods = map[string]struct{}{
	"Info":  {},
	"Warn":  {},
	"Error": {},
	"Debug": {},
}

func runWithConfig(cfg *config.Config) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		detectors := []*detector.PatternDetector{}
		if cfg.Rules.CustomPatterns {
			for _, customPattern := range cfg.CustomPatterns {
				pattern := customPattern.Pattern
				detector := detector.NewPatternDetector(pattern)
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
					isLog = isLogger(cfg.KnownPacks, obj)
				case *types.Var:
					isLog = isLoggerVar(cfg.KnownPacks, obj.Type())
				case *types.Func:
					isLog = isLoggerFunc(cfg.KnownPacks, obj)
				}

				if !isLog {
					return true
				}

				if !isLogMethod(sel) {
					return true
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

// Check if the call expression is a log call
func isLogger(knownPacks map[string]config.PackInfo, pkgName *types.PkgName) bool {
	pkgPath := pkgName.Imported().Path()

	if knownPack, ok := knownPacks[pkgPath]; !ok || !knownPack.Enabled {
		return false
	}
	return true
}

// Check if the call expression is a logger variable
func isLoggerVar(knownPacks map[string]config.PackInfo, t types.Type) bool {
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
	pkgPath := pkg.Path()
	if _, ok := knownPacks[pkgPath]; ok {
		return true
	}
	return false
}

// Check if the call expression is a function that returns a logger variable
func isLoggerFunc(knownPacks map[string]config.PackInfo, fn *types.Func) bool {
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}

	res := sig.Results()
	for i := range res.Len() {
		if isLoggerVar(knownPacks, res.At(i).Type()) {
			return true
		}
	}
	return false
}

// Check if its method is a log method
func isLogMethod(sel *ast.SelectorExpr) bool {
	methodName := sel.Sel.Name
	_, ok := LogMethods[methodName]
	return ok
}
