package tokenizer

import (
	"strings"

	"github.com/CHIRANTAN-001/zeno/internal/stopwords"
)

func Tokenize(text string, filter *stopwords.Filter) []string {
	words := strings.Fields((strings.ToLower(text)))

	var result []string

	for _, word := range words {
		if !filter.IsStopWord(word) {
			result = append(result, word)
		}
	}

	return result
}
