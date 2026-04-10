package tokenizer

import (
	"github.com/CHIRANTAN-001/zeno/internal/stopwords"
	"regexp"
	"strings"
)

// Non-alphanumeric characters
var nonAlphanumericChars = regexp.MustCompile("[^a-z0-9]+")

func Tokenize(text string, filter *stopwords.Filter) []string {
	tokens := strings.Fields((strings.ToLower(text)))
	var filteredTokens []string

	for _, token := range tokens {
		// Remove non-alphanumeric characters
		token = nonAlphanumericChars.ReplaceAllString(token, "")
		if token != "" && !filter.IsStopWord(token) {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return filteredTokens
}
