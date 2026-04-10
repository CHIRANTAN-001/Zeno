package index

type InvertedIndex struct {
	Index map[string]map[int]int
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index: make(map[string]map[int]int),
	}
}

func (idx *InvertedIndex) Add(docID int, tokens []string) {
	for _, token := range tokens {
		if _, exists := idx.Index[token]; !exists {
			idx.Index[token] = make(map[int]int)
		}
		idx.Index[token][docID] = idx.Index[token][docID] + 1
	}
}
