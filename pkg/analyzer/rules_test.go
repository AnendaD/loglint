package analyzer

import (
	"testing"
)

// isEnglishLetter

func TestIsEnglishLetter(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'м', false}, // кириллица
		{'é', false}, // латиница с диакритикой
		{'中', false}, // китайский
		{'1', false}, // цифра — не буква вообще, но функция вернёт false
		{' ', false},
	}
	for _, tc := range tests {
		got := isEnglishLetter(tc.r)
		if got != tc.want {
			t.Errorf("isEnglishLetter(%q) = %v, want %v", tc.r, got, tc.want)
		}
	}
}

// isSpecialOrEmoji

func TestIsSpecialOrEmoji(t *testing.T) {
	tests := []struct {
		r           rune
		wantSpecial bool // первое возвращаемое значение
	}{
		{'a', false},
		{'Z', false},
		{'1', false},
		{' ', false}, // пробел — не специальный
		{'-', false}, // разрешённый символ
		{'_', false},
		{'.', false},
		{',', false},
		{':', false},
		{'/', false},
		{'\\', false},
		{'!', true}, // запрещён
		{'?', true},
		{'🚀', true}, // эмодзи
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

// hasSpecialOrEmoji

func TestHasSpecialOrEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"server started", false},
		{"connection failed", false},
		{"something went wrong", false},
		{"server started!", true},                  // восклицательный знак
		{"server started🚀", true},                  // эмодзи
		{"connection failed!!!", true},             // повторяющиеся !
		{"warning: something went wrong...", true}, // многоточие
		{"failed to connect to db", false},
		{"user-agent received", false}, // дефис разрешён
		{"key_value pair", false},      // подчёркивание разрешено
		{"port 8080/tcp", false},       // слэш разрешён
		{"path: /usr/bin", false},
	}
	for _, tc := range tests {
		got := hasSpecialOrEmoji(tc.input)
		if got != tc.want {
			t.Errorf("hasSpecialOrEmoji(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// removeSpecialOrEmoji

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
		{"key: value", "key: value"}, // одно двоеточие разрешено
	}
	for _, tc := range tests {
		got := removeSpecialOrEmoji(tc.input)
		if got != tc.want {
			t.Errorf("removeSpecialOrEmoji(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// redactKeyWord

func TestRedactKeyWord(t *testing.T) {
	tests := []struct {
		s       string
		keyword string
		want    string
	}{
		{"user password: abc", "password", "user [REDACTED]: abc"},
		{"api_key=123", "api_key", "[REDACTED]=123"},
		{"no sensitive info here", "password", "no sensitive info here"},
		{"PASSWORD uppercase", "password", "[REDACTED] uppercase"},       // case-insensitive
		{"token: abc token: def", "token", "[REDACTED]: abc token: def"}, // только первое вхождение
		{"", "password", ""},
	}
	for _, tc := range tests {
		got := redactKeyWord(tc.s, tc.keyword)
		if got != tc.want {
			t.Errorf("redactKeyWord(%q, %q) = %q, want %q", tc.s, tc.keyword, got, tc.want)
		}
	}
}
