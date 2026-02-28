package detector_test

import (
	"selectellinter/pkg/analyzer/detector"
	"testing"
)

func TestAWSDetector_Detect(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    int // количество совпадений
	}{
		{
			name:    "aws key detected",
			pattern: `(?i)AKIA[0-9A-Z]{16}`,
			input:   "key: AKIAIOSFODNN7EXAMPLE",
			want:    1,
		},
		{
			name:    "aws key not present",
			pattern: `(?i)AKIA[0-9A-Z]{16}`,
			input:   "just a regular string",
			want:    0,
		},
		{
			name:    "jwt detected",
			pattern: `eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`,
			input:   "token: eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyIn0.abc123",
			want:    1,
		},
		{
			name:    "multiple matches",
			pattern: `\d+`,
			input:   "ports 8080 and 9090",
			want:    2,
		},
		{
			name:    "empty string",
			pattern: `(?i)AKIA[0-9A-Z]{16}`,
			input:   "",
			want:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := detector.NewAWSDetector(tc.pattern)
			got := d.Detect(tc.input)
			if len(got) != tc.want {
				t.Errorf("Detect(%q) = %d matches, want %d", tc.input, len(got), tc.want)
			}
		})
	}
}

func TestAWSDetector_Replace(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		input       string
		replacement string
		want        string
	}{
		{
			name:        "replace aws key",
			pattern:     `(?i)AKIA[0-9A-Z]{16}`,
			input:       "key AKIAIOSFODNN7EXAMPLE here",
			replacement: "[REDACTED]",
			want:        "key [REDACTED] here",
		},
		{
			name:        "no match — string unchanged",
			pattern:     `(?i)AKIA[0-9A-Z]{16}`,
			input:       "nothing to replace",
			replacement: "[REDACTED]",
			want:        "nothing to replace",
		},
		{
			name:        "replace all occurrences",
			pattern:     `\d+`,
			input:       "port 8080 and 9090",
			replacement: "XXXX",
			want:        "port XXXX and XXXX",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := detector.NewAWSDetector(tc.pattern)
			got := d.Replace(tc.input, tc.replacement)
			if got != tc.want {
				t.Errorf("Replace(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
