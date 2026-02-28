package analyzer

import (
	"selectellinter/config"

	"golang.org/x/tools/go/analysis"
)

// Returns new analyzer
func NewAnalyzer(cfg *config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "linter",
		Doc:  "Linter analyzer demo",
		Run:  runWithConfig(cfg),
	}
}
