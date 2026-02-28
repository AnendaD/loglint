package main

import (
	"selectellinter/config"
	"selectellinter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := config.MustLoad()
	singlechecker.Main(analyzer.NewAnalyzer(cfg))

}
