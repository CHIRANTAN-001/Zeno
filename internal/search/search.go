package search

import (
	"github.com/CHIRANTAN-001/zeno/internal/index"
)

func Search(idx *index.InvertedIndex, tokens []string) []int {
	var result map[int]int

	for i, token := range tokens {
		docs := idx.Index[token]
		if docs == nil {
			return []int{}
		}
		if i == 0 {
			result = docs
		} else {
			temp := make(map[int]int)
			for docID := range result {
				if docs[docID] > 0 {
					temp[docID] = docs[docID]
				}
			}
			result = temp
		}
	}

	var output []int
	for docID := range result {
		output = append(output, docID)
	}
	return output
}
