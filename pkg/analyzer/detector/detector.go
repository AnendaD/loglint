package detector

import "regexp"

type AWSDetector struct {
	regexp *regexp.Regexp
}

// Returns new AWSDetector
func NewAWSDetector(r string) *AWSDetector {
	return &AWSDetector{
		regexp: regexp.MustCompile(r),
	}
}

// Detect returns all matches in the string
func (r *AWSDetector) Detect(s string) []string {
	return r.regexp.FindAllString(s, -1)
}

// Replace replaces all matches in the string
func (r *AWSDetector) Replace(s, rep string) string {
	return r.regexp.ReplaceAllString(s, rep)
}
