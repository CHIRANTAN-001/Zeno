package main

import (
	"fmt"
	"sort"

	"github.com/CHIRANTAN-001/zeno/internal/index"
	"github.com/CHIRANTAN-001/zeno/internal/search"
	"github.com/CHIRANTAN-001/zeno/internal/stopwords"
	"github.com/CHIRANTAN-001/zeno/internal/tokenizer"
)

type searchHit struct {
	ID   int
	Text string
}

func main() {
	// DOCUMENTS
	docs := map[int]string{
		1:  "The quick brown fox jumps over the lazy dog",
		2:  "The dog barked at the cat",
		3:  "The cat, cat chased the mouse",
		4:  "The mouse ran away from the cat",
		5:  "The cat caught the mouse",
		6:  "The mouse escaped from the cat",
		7:  "The cat chased the mouse",
		8:  "The mouse ran away from the cat",
		9:  "The cat caught the mouse",
		10: "The mouse escaped from the cat",
	}

	// INIT
	idx := index.NewInvertedIndex()
	filter := stopwords.NewFilter()

	// BUILD INDEX
	for id, text := range docs {
		tokens := tokenizer.Tokenize(text, filter)
		idx.Add(id, tokens)
	}

	fmt.Println("index", idx.Index)

	// QUERY
	query := "cat mouse"
	queryTokens := tokenizer.Tokenize(query, filter)

	results := search.Search(idx, queryTokens)

	// BUILD HITS
	hits := make([]searchHit, 0, len(results))
	for _, docID := range results {
		hits = append(hits, searchHit{
			ID:   docID,
			Text: docs[docID],
		})
	}

	// Sort results for consistent output
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].ID < hits[j].ID
	})

	// PRINT
	printResults(query, hits)
	printIndex(idx.Index)
}

// DISPLAY HELPERS

func printResults(query string, hits []searchHit) {
	fmt.Println("\n================ SEARCH RESULTS ================")

	fmt.Printf("\nQuery: %q\n", query)
	fmt.Printf("\nMatched Documents (%d):\n\n", len(hits))

	if len(hits) == 0 {
		fmt.Println("No results found.")
	} else {
		for _, h := range hits {
			fmt.Printf("[%d] %s\n", h.ID, h.Text)
		}
	}

	fmt.Println("\n================================================")
}

func printIndex(idx map[string]map[int]int) {
	fmt.Println("\n================ INVERTED INDEX ================\n")

	// Sort terms for deterministic output
	var terms []string
	for term := range idx {
		terms = append(terms, term)
	}
	sort.Strings(terms)

	for _, term := range terms {
		docIDs := make([]int, 0, len(idx[term]))
		for id := range idx[term] {
			docIDs = append(docIDs, id)
		}
		sort.Ints(docIDs)
		fmt.Printf("%-12s -> %v\n", term, docIDs)
	}

	fmt.Println("\n================================================")
}