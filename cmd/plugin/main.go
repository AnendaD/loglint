package main

import (
	"selectellinter/config"
	"selectellinter/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	cfg := config.MustLoad()
	return []*analysis.Analyzer{analyzer.NewAnalyzer(cfg)}, nil
}