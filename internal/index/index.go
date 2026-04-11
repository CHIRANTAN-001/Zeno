package index

import (
	"encoding/gob"
	"os"
)

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

func (idx *InvertedIndex) SaveToDisk(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(idx.Index)
}

func (idx *InvertedIndex) LoadFromDisk(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return  err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&idx.Index)
}