package main

import (
	"linter/config"
	"linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := config.MustLoad()
	singlechecker.Main(analyzer.NewAnalyzer(cfg))
}
