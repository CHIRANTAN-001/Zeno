package index

type InvertedIndex struct {
	Index map[string]map[int]bool
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index: make(map[string]map[int]bool),
	}
}

func (idx *InvertedIndex) Add(docID int, tokens []string) {
	for _, token := range tokens {
		if _, exists := idx.Index[token]; !exists {
			idx.Index[token] = make(map[int]bool)
		}
		idx.Index[token][docID] = true
	}
}