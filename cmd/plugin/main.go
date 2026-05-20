package main

import (
	"linter/config"
	"linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

type PluginSymbol struct {
	Analyzer *analysis.Analyzer
}

func (p *PluginSymbol) Build() *analysis.Analyzer {
	cfg := config.MustLoad()
	return analyzer.NewAnalyzer(cfg)
}

var Plugin PluginSymbol
