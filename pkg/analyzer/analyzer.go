package analyzer

import (
	"selectellinter/config"

	"golang.org/x/tools/go/analysis"
)

// Returns new analyzer
func NewAnalyzer(cfg *config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "loglinter",
		Doc:  "Go linter for log message validation",
		Run:  runWithConfig(cfg),
	}
}
