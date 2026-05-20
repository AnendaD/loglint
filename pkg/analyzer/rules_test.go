package analyzer

import (
	"testing"
)

func TestIsEnglishLetter(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'м', false},
		{'é', false},
		{'中', false},
		{'1', false},
		{' ', false},
	}
	for _, tc := range tests {
		got := isEnglishLetter(tc.r)
		if got != tc.want {
			t.Errorf("isEnglishLetter(%q) = %v, want %v", tc.r, got, tc.want)
		}
	}
}

func TestIsSpecialOrEmoji(t *testing.T) {
	tests := []struct {
		r           rune
		wantSpecial bool
	}{
		{'a', false},
		{'Z', false},
		{'1', false},
		{' ', false},
		{'-', false},
		{'_', false},
		{'.', false},
		{',', false},
		{':', false},
		{'/', false},
		{'\\', false},
		{'!', true},
		{'?', true},
		{'🚀', true},
		{'@', true},
		{'#', true},
		{'(', true},
		{')', true},
	}
	for _, tc := range tests {
		gotSpecial, _ := isSpecialOrEmoji(tc.r)
		if gotSpecial != tc.wantSpecial {
			t.Errorf("isSpecialOrEmoji(%q) special = %v, want %v", tc.r, gotSpecial, tc.wantSpecial)
		}
	}
}

func TestHasSpecialOrEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"server started", false},
		{"connection failed", false},
		{"something went wrong", false},
		{"server started!", true},
		{"server started🚀", true},
		{"connection failed!!!", true},
		{"warning: something went wrong...", true},
		{"failed to connect to db", false},
		{"user-agent received", false},
		{"key_value pair", false},
		{"port 8080/tcp", false},
		{"path: /usr/bin", false},
	}
	for _, tc := range tests {
		got := hasSpecialOrEmoji(tc.input)
		if got != tc.want {
			t.Errorf("hasSpecialOrEmoji(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestRemoveSpecialOrEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"server started!", "server started"},
		{"connection failed!!!", "connection failed"},
		{"server started🚀", "server started"},
		{"something went wrong...", "something went wrong"},
		{"hello world", "hello world"},
		{"key: value", "key: value"},
	}
	for _, tc := range tests {
		got := removeSpecialOrEmoji(tc.input)
		if got != tc.want {
			t.Errorf("removeSpecialOrEmoji(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
