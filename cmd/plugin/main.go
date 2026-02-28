package main

import (
	"github.com/AnendaD/loglint/config"
	"github.com/AnendaD/loglint/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	cfg := config.MustLoad()
	return []*analysis.Analyzer{analyzer.NewAnalyzer(cfg)}, nil
}