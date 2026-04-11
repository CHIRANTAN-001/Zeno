package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CHIRANTAN-001/zeno/internal/index"
	"github.com/CHIRANTAN-001/zeno/internal/parser"
	"github.com/CHIRANTAN-001/zeno/internal/search"
	"github.com/CHIRANTAN-001/zeno/internal/stopwords"
	"github.com/CHIRANTAN-001/zeno/internal/tokenizer"
)

var (
	indexPath = "zeno_bin"
	docsPath  = "zeno_docs"
)
// wikipedia file path
const wikiDump = "./data/simplewiki-latest-pages-articles.xml.bz2"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ")
		fmt.Println("  go run main.go index")
		fmt.Println("  go run main.go search <query>")
		return
	}

	switch os.Args[1] {
	case "index":
		runIndexer()
	// case "search":
	// 	if len(os.Args) < 1 {
	// 		fmt.Println("Usage: go run main.go search <query>")
	// 		return
	// 	}
	// 	query := strings.Join(os.Args[2:], " ")
	// 	runSearch(query)
	case "serve":
		runServer()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func runIndexer() {
	fmt.Println("Starting indexer...")

	idx := index.NewInvertedIndex()
	filter := stopwords.NewFilter()
	docs := index.NewDocStore()

	articles, errc := parser.StreamArticles(wikiDump)
	docID := 0

	for article := range articles {
		docID++
		docs.Add(docID, article.Title, article.Body)
		fullText := article.Title + " " + article.Body
		tokens := tokenizer.Tokenize(fullText, filter)
		idx.Add(docID, tokens)

		if docID%10_000 == 0 {
			fmt.Println("Indexed", docID, "articles")
		}
	}

	if err := <-errc; err != nil {
		fmt.Println("Parser error:", err)
	}

	fmt.Printf("Total articles indexed: %d\n", docID)

	// save both index and docs to disk
	if err := idx.SaveToDisk(indexPath); err != nil {
		fmt.Println("Failed to save index:", err)
		return
	}
	fmt.Println("Index saved")

	if err := docs.SaveToDisk(docsPath); err != nil {
		fmt.Println("Failed to save docs:", err)
		return
	}

	fmt.Println("Done. Index saved to", indexPath)
}

// func runSearch(query string) {
// 	t0 := time.Now()

// 	idx := index.NewInvertedIndex()
// 	idx.LoadFromDisk(indexPath)
// 	fmt.Printf("Index loaded in %v\n", time.Since(t0))

// 	t1 := time.Now()
// 	docs := index.NewDocStore()
// 	docs.LoadFromDisk(docsPath)
// 	fmt.Printf("Docs loaded in %v\n", time.Since(t1))

// 	t2 := time.Now()
// 	filter := stopwords.NewFilter()
// 	queryTokens := tokenizer.Tokenize(query, filter)
// 	results := search.Search(idx, queryTokens)
// 	fmt.Printf("Search took %v\n", time.Since(t2))

// 	fmt.Printf("\nResults for %q:\n\n", query)
// 	for _, id := range results {
// 		fmt.Printf("[%d] %s\n", id, docs.Get(id))
// 	}
// 	fmt.Printf("\n%d results found\n", len(results))
// }

func runServer() {
	fmt.Println("Loading index...")
	t0 := time.Now()

	idx := index.NewInvertedIndex()
	if err := idx.LoadFromDisk(indexPath); err != nil {
		fmt.Println("Failed to load index:", err)
		return
	}

	docs := index.NewDocStore()
	if err := docs.LoadFromDisk(docsPath); err != nil {
		fmt.Println("Docs not found. Run: go run main.go index")
		return
	}

	fmt.Printf("Docs entries loaded: %d\n", len(docs.Docs))

	fmt.Printf("Ready in %v — listening on http://localhost:8080\n", time.Since(t0))

	filter := stopwords.NewFilter()

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "missing ?q=", 400)
			return
		}

		t := time.Now()
		tokens := tokenizer.Tokenize(query, filter)
		results := search.Search(idx, tokens)

		// build response with titles
		type Hit struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
		}
		hits := make([]Hit, 0, len(results))
		for _, id := range results {
			doc := docs.Get(id)
			snippet := doc.Body
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}
			hits = append(hits, Hit{
				ID:      id,
				Title:   doc.Title,
				Snippet: snippet,
			})
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]any{
			"query":   query,
			"count":   len(results),
			"results": hits,
			"took_ms": time.Since(t).Milliseconds(),
		})
	})
	http.ListenAndServe(":8080", nil)
}
