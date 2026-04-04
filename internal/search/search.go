package search

import "github.com/CHIRANTAN-001/zeno/internal/index"

func Search(idx *index.InvertedIndex, tokens []string) []int {
	var result map[int]bool

	for i, token := range tokens {
		docs, exists := idx.Index[token]
		if !exists {
			return []int{}
		}

		if i == 0 {
			result = docs
		} else {
			temp := make(map[int]bool)
			for docID := range result {
				if docs[docID] {
					temp[docID] = true
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
