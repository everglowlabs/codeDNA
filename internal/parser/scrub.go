package parser

import (
	"regexp"
	"strings"
)

var (
	// Matches double-quoted strings (handling escaped quotes)
	doubleQuoteRegex = regexp.MustCompile(`"([^"\\]|\\.)*"`)
	// Matches single-quoted strings
	singleQuoteRegex = regexp.MustCompile(`'([^'\\]|\\.)*'`)
	// Matches backtick strings (JS/Go/Rust backticks)
	backtickRegex = regexp.MustCompile("`([^`\\\\]|\\\\.)*`")
	// Matches Python triple quotes (double and single) using RE2 compatible non-greedy dotall
	pyTripleDoubleRegex = regexp.MustCompile(`(?s)""".*?"""`)
	pyTripleSingleRegex = regexp.MustCompile(`(?s)'''.*?'''`)
)

// ScrubStrings removes all string literals from the code content, replacing them with [scrubbed]
func ScrubStrings(content string) string {
	content = pyTripleDoubleRegex.ReplaceAllString(content, `__TRIPLE_DOUBLE_SCRUBBED__`)
	content = pyTripleSingleRegex.ReplaceAllString(content, `__TRIPLE_SINGLE_SCRUBBED__`)
	content = doubleQuoteRegex.ReplaceAllString(content, `"[scrubbed]"`)
	content = singleQuoteRegex.ReplaceAllString(content, `'[scrubbed]'`)
	content = backtickRegex.ReplaceAllString(content, "`[scrubbed]`")
	content = strings.ReplaceAll(content, "__TRIPLE_DOUBLE_SCRUBBED__", `"""[scrubbed]"""`)
	content = strings.ReplaceAll(content, "__TRIPLE_SINGLE_SCRUBBED__", `'''[scrubbed]'''`)
	return content
}
