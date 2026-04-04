package stopwords

import "strings"

type Filter struct {
	words map[string]struct{}
}

func NewFilter() *Filter {
	return &Filter{
		words: map[string]struct{}{
			"a":          {},
			"an":         {},
			"and":        {},
			"are":        {},
			"as":         {},
			"at":         {},
			"by":         {},
			"for":        {},
			"from":       {},
			"in":         {},
			"is":         {},
			"of":         {},
			"on":         {},
			"or":         {},
			"out":        {},
			"over":       {},
			"up":         {},
			"with":       {},
			"without":    {},
			"the":        {},
			"to":         {},
			"was":        {},
			"we":         {},
			"were":       {},
			"what":       {},
			"when":       {},
			"where":      {},
			"why":        {},
			"you":        {},
			"your":       {},
			"yours":      {},
			"yourself":   {},
			"yourselves": {},
			"themselves": {},
		},
	}
}

func (f *Filter) IsStopWord(word string) bool {
	_, exists := f.words[strings.ToLower(word)]
	return exists
}

func (f *Filter) Add(word string) {
	f.words[word] = struct{}{}
}

func (f *Filter) Remvoe(word string) {
	delete(f.words, word)
}