package tokenizer

import (
	"github.com/CHIRANTAN-001/zeno/internal/stopwords"
	"regexp"
	"strings"
)

var nonAlpha = regexp.MustCompile("[^a-z0-9]+")

func Tokenize(text string, filter *stopwords.Filter) []string {
	words := strings.Fields((strings.ToLower(text)))
	var result []string

	for _, word := range words {
		word = nonAlpha.ReplaceAllString(word, "")
		if word != "" && !filter.IsStopWord(word) {
			result = append(result, word)
		}
	}

	return result
}
