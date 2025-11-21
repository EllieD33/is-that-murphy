package utils_test

import (
	"testing"

	"github.com/ellied33/is-that-murphy/utils"
)

func TestCanonical(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple lowercase", "hello", "hello"},
		{"leading spaces", "   hello", "hello"},
		{"trailing spaces", "hello   ", "hello"},
		{"both sides spaces", "   hello   ", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"mixed case", "HeLlO", "hello"},
		{"unicode uppercase", "ÅNGSTRÖM", "ångström"},
		{"newline trimmed", "\nhello\n", "hello"},
		{"tab trimmed", "\thello\t", "hello"},
		{"carriage return trimmed", "\rhello\r", "hello"},
		{"crlf", "\r\nhello\r\n", "hello"},
		{"empty string", "", ""},
		{"just spaces", "     ", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := utils.Canonical(tc.input)
			if got != tc.want {
				t.Fatalf("Canonical(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
