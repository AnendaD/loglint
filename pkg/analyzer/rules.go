package analyzer

import (
	"go/ast"
	"go/token"
	"selectellinter/config"
	"selectellinter/pkg/analyzer/detector"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// Log message validation rules:
//   - Must start with a lowercase letter
//   - Must be in English only
//   - No emojis or special characters allowed
//   - No sensitive information allowed
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr, cfg *config.Config, detectors []*detector.AWSDetector) {
	if len(call.Args) == 0 {
		return
	}
	// - Must start with a lowercase letter
	msgArg := call.Args[0]
	rules := cfg.Rules
	if rules.Lowercase {
		if lit, ok := msgArg.(*ast.BasicLit); ok {
			raw := lit.Value
			if len(raw) > 2 {
				fChar, size := utf8.DecodeRuneInString(raw[1:])
				if unicode.IsUpper(fChar) {
					fixed := raw[:1] + strings.ToLower(string(fChar)) + raw[1+size:]
					d := analysis.Diagnostic{
						Pos:     lit.Pos(),
						Message: "Log must start with a lowercase letter",
					}
					if cfg.AutoFix.Enabled && cfg.AutoFix.Lowercase {
						d.SuggestedFixes = []analysis.SuggestedFix{{
							Message: "Convert first letter to lowercase",
							TextEdits: []analysis.TextEdit{{
								Pos:     lit.Pos(),
								End:     lit.End(),
								NewText: []byte(fixed),
							}},
						}}
					}
					pass.Report(d)
				}
			}
		}
	}
	var literals []*ast.BasicLit
	var identNames []string
	for _, arg := range call.Args {
		collectLiterals(arg, &literals)
		collectIdents(arg, &identNames)
	}

	//eng letters
	if cfg.Rules.English {
		flag := false
		for _, lit := range literals {
			for _, ch := range lit.Value[1 : len(lit.Value)-1] {
				if unicode.IsLetter(ch) && !isEnglishLetter(ch) {
					pass.Reportf(call.Pos(), "log message must be in English (found: %c)", ch)
					flag = true
					break
				}
			}
			if flag {
				break
			}
		}
	}

	//special chars or emojis
	if cfg.Rules.SpecialChars {
		for _, lit := range literals {
			raw := lit.Value[1 : len(lit.Value)-1]
			if hasSpecialOrEmoji(raw) {
				fixed := `"` + removeSpecialOrEmoji(raw) + `"`
				d := analysis.Diagnostic{
					Pos:     lit.Pos(),
					Message: "log message must not contain special characters or emoji",
				}
				if cfg.AutoFix.Enabled && cfg.AutoFix.SpecialChars {
					d.SuggestedFixes = []analysis.SuggestedFix{{
						Message: "Remove special characters and emoji",
						TextEdits: []analysis.TextEdit{{
							Pos:     lit.Pos(),
							End:     lit.End(),
							NewText: []byte(fixed),
						}},
					}}
				}
				pass.Report(d)
			}
		}
	}

	//keywords
	if rules.SensitiveKeywords {
		for _, ident := range identNames {
			lower := strings.ToLower(ident)
			for _, kw := range cfg.SensitiveKeywords {
				if strings.Contains(lower, strings.ToLower(kw)) {
					pass.Reportf(call.Pos(), "log message contains sensitive keyword: %q", kw)
					if cfg.AutoFix.Enabled && cfg.AutoFix.SensitiveKeywords {
						break
					}
				}
			}
		}
	}

	//patterns
	if cfg.Rules.CustomPatterns {
		for i, lit := range literals {
			raw := lit.Value[1 : len(lit.Value)-1]
			for j, det := range detectors {
				if matches := det.Detect(raw); len(matches) == 0 {
					continue
				}
				d := analysis.Diagnostic{
					Pos:     lit.Pos(),
					Message: cfg.CustomPatterns[j].Message,
				}
				if cfg.AutoFix.Enabled && cfg.AutoFix.CustomPatterns {
					fixed := det.Replace(raw, cfg.CustomPatterns[j].Replacement)
					d.SuggestedFixes = []analysis.SuggestedFix{{
						Message: cfg.CustomPatterns[j].Message,
						TextEdits: []analysis.TextEdit{{
							Pos:     literals[i].Pos(),
							End:     literals[i].End(),
							NewText: []byte(`"` + fixed + `"`),
						}},
					}}
				}
				pass.Report(d)
			}
		}
	}
}

func collectLiterals(expr ast.Expr, res *[]*ast.BasicLit) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if len(e.Value) >= 2 {
				*res = append(*res, e)
			}
		}
	case *ast.BinaryExpr:
		collectLiterals(e.X, res)
		collectLiterals(e.Y, res)
	case *ast.CallExpr:
		for _, arg := range e.Args {
			collectLiterals(arg, res)
		}
	}
}

func collectIdents(expr ast.Expr, res *[]string) {
	switch e := expr.(type) {
	case *ast.Ident:
		*res = append(*res, e.Name)
	case *ast.BinaryExpr:
		collectIdents(e.X, res)
		collectIdents(e.Y, res)
	case *ast.CallExpr:
		collectIdents(e.Fun, res)
		for _, arg := range e.Args {
			collectIdents(arg, res)
		}
	case *ast.SelectorExpr:
		collectIdents(e.X, res)
		*res = append(*res, e.Sel.Name)
	case *ast.IndexExpr:
		collectIdents(e.X, res)
		collectIdents(e.Index, res)
	case *ast.StarExpr:
		collectIdents(e.X, res)
	case *ast.UnaryExpr:
		collectIdents(e.X, res)
	case *ast.TypeAssertExpr:
		collectIdents(e.X, res)
	case *ast.SliceExpr:
		collectIdents(e.X, res)
		if e.Low != nil {
			collectIdents(e.Low, res)
		}
		if e.High != nil {
			collectIdents(e.High, res)
		}
	}
}

func isEnglishLetter(r rune) bool {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return true
	}
	return false
}

func isSpecialOrEmoji(r rune) (bool, bool) {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
		return false, false
	}
	allowed := " -_.,:/\\"
	if strings.ContainsRune(allowed, r) {
		return false, true
	}
	return true, false
}

func hasSpecialOrEmoji(s string) bool {
	isPrevAllowed := false
	for _, r := range s {
		f, isAllowed := isSpecialOrEmoji(r)
		if f || (isAllowed && isPrevAllowed) {
			return true
		}
		isPrevAllowed = isAllowed
	}
	return false
}

func removeSpecialOrEmoji(s string) string {
	res := []rune{}
	isPrevAllowed := false
	duplicate := false
	for _, r := range s {
		isSpec, isAllowed := isSpecialOrEmoji(r)
		if isPrevAllowed && isAllowed {
			duplicate = true
		}
		if !isSpec && !(isPrevAllowed && isAllowed) {
			if duplicate {
				res = res[:len(res)-1]
				duplicate = false
			}
			res = append(res, r)
		}
		isPrevAllowed = isAllowed
	}
	if duplicate {
		res = res[:len(res)-1]
	}
	return string(res)
}
