package search

import (
	"github.com/CHIRANTAN-001/zeno/internal/index"
)

func Search(idx *index.InvertedIndex, tokens []string) []int {
	var result map[int]int

	for i, token := range tokens {
		docs, exists := idx.Index[token]
		if !exists {
			return []int{}
		}
		
		if i == 0 {
			result = docs
		} else {
			temp := make(map[int]int)
			for docID := range result {
				if _, found := docs[docID]; found {
					temp[docID] = docs[docID]
				}
			}
			result = temp
		}
	}

	output := make([]int, 0, len(result))
	for docID := range result {
		output = append(output, docID)
	}
	return output
}
