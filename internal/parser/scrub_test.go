package parser

import "testing"

func TestScrubStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Double quotes",
			input:    `x := "hello world"`,
			expected: `x := "[scrubbed]"`,
		},
		{
			name:     "Single quotes",
			input:    `char := 'a'`,
			expected: `char := '[scrubbed]'`,
		},
		{
			name:     "Backticks",
			input:    "s := `hello\nworld`",
			expected: "s := `[scrubbed]`",
		},
		{
			name:     "Python triple double quotes",
			input:    `"""docstring\nhere"""`,
			expected: `"""[scrubbed]"""`,
		},
		{
			name:     "Python triple single quotes",
			input:    `'''docstring\nhere'''`,
			expected: `'''[scrubbed]'''`,
		},
		{
			name:     "Mixed strings",
			input:    `func test() { x := "one"; y := 'b'; z := ` + "`" + `three` + "`" + ` }`,
			expected: `func test() { x := "[scrubbed]"; y := '[scrubbed]'; z := ` + "`" + `[scrubbed]` + "`" + ` }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScrubStrings(tt.input)
			if got != tt.expected {
				t.Errorf("ScrubStrings() = %q, want %q", got, tt.expected)
			}
		})
	}
}
