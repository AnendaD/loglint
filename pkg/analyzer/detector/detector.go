package detector

import "regexp"

type PatternDetector struct {
	regexp *regexp.Regexp
}

// Returns new PatternDetector
func NewPatternDetector(r string) *PatternDetector {
	return &PatternDetector{
		regexp: regexp.MustCompile(r),
	}
}

// Detect returns all matches in the string
func (r *PatternDetector) Detect(s string) []string {
	return r.regexp.FindAllString(s, -1)
}

// Replace replaces all matches in the string
func (r *PatternDetector) Replace(s, rep string) string {
	return r.regexp.ReplaceAllString(s, rep)
}
