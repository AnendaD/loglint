package main

import (
	"github.com/AnendaD/loglint/config"
	"github.com/AnendaD/loglint/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := config.MustLoad()
	singlechecker.Main(analyzer.NewAnalyzer(cfg))

}
