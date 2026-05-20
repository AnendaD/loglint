package analyzer_test

import (
	"testing"

	"linter/config"
	"linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

var testdataPath = analysistest.TestData()

func TestAnalyzer_LowercaseRule(t *testing.T) {
	cfg := &config.Config{
		KnownPacks: map[string]config.PackInfo{
			"log/slog": {
				Version: "1.0.0",
				Enabled: true,
			},
			"go.uber.org/zap": {
				Version: "1.26.0",
				Enabled: true,
			},
		},
		Rules: config.Rules{
			Lowercase: true,
		},
		SensitiveKeywords: []string{},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataPath, a, "lowercase")
}

func TestAnalyzer_EnglishRule(t *testing.T) {
	cfg := &config.Config{
		KnownPacks: map[string]config.PackInfo{
			"log/slog": {
				Version: "1.0.0",
				Enabled: true,
			},
			"go.uber.org/zap": {
				Version: "1.26.0",
				Enabled: true,
			},
		},
		Rules: config.Rules{
			English: true,
		},
		SensitiveKeywords: []string{},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataPath, a, "english")
}

func TestAnalyzer_SpecialCharsRule(t *testing.T) {
	cfg := &config.Config{
		KnownPacks: map[string]config.PackInfo{
			"log/slog": {
				Version: "1.0.0",
				Enabled: true,
			},
			"go.uber.org/zap": {
				Version: "1.26.0",
				Enabled: true,
			},
		},
		Rules: config.Rules{
			SpecialChars: true,
		},
		SensitiveKeywords: []string{},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataPath, a, "specialchars")
}

func TestAnalyzer_SensitiveKeywordsRule(t *testing.T) {
	cfg := &config.Config{
		KnownPacks: map[string]config.PackInfo{
			"log/slog": {
				Version: "1.0.0",
				Enabled: true,
			},
			"go.uber.org/zap": {
				Version: "1.26.0",
				Enabled: true,
			},
		},
		Rules: config.Rules{
			SensitiveKeywords: true,
		},
		SensitiveKeywords: []string{"password", "token", "apiKey", "secret"},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataPath, a, "sensitive")
}

func TestAnalyzer_CustomPatternsRule(t *testing.T) {
	cfg := &config.Config{
		KnownPacks: map[string]config.PackInfo{
			"log/slog": {
				Version: "1.0.0",
				Enabled: true,
			},
			"go.uber.org/zap": {
				Version: "1.26.0",
				Enabled: true,
			},
		},
		Rules: config.Rules{
			CustomPatterns: true,
		},
		SensitiveKeywords: []string{},
		CustomPatterns: []config.CustomPattern{
			{
				Name:    "aws_key",
				Pattern: `(?i)AKIA[0-9A-Z]{16}`,
				Message: "AWS access key detected",
			},
		},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataPath, a, "custom")
}
