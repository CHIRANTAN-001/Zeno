package parser

import (
	"compress/bzip2"
	"encoding/xml"
	"io"
	"os"
)

type Article struct {
	Title string
	Body  string
}

type WikiPage struct {
	Title    string `xml:"title"`
	Redirect struct {
		Title string `xml:"title"`
	} `xml:"redirect"`
	Revision struct {
		Text string `xml:"text"`
	} `xml:"revision"`
}

// StreamArticles opens the Wikipedia bz2 XML dump and streams
// Articles one by one into the returned channel
func StreamArticles(path string) (<-chan Article, <-chan error) {
	articles := make(chan Article, 100)
	errc := make(chan error, 1)

	go func() {
		defer close(articles)
		defer close(errc)

		f, err := os.Open(path)
		if err != nil {
			errc <- err
			return
		}
		defer f.Close()

		// decompress zip file on the go 
		br := bzip2.NewReader(f)
		decoder := xml.NewDecoder(br)

		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}

			// look for <page> tag
			if se, ok := token.(xml.StartElement); ok && se.Name.Local == "page" {
				var page WikiPage
				if err := decoder.DecodeElement(&page, &se); err != nil {
					continue
				}

				// skip redirect pages
				if page.Redirect.Title != "" {
					continue
				}

				// skip empty articles
				if page.Revision.Text == "" {
					continue
				}

				articles <- Article{
					Title: page.Title,
					Body:  page.Revision.Text,
				}
			}
		}
	}()

	return articles, errc
}
