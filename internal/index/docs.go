package index

import (
	"encoding/gob"
	"os"
)

type Doc struct {
    Title string
    Body  string
}

type DocStore struct {
	Docs map[int]Doc
}

func NewDocStore() *DocStore {
	return &DocStore{
		Docs: make(map[int]Doc),
	}
}

func (d *DocStore) Add(docID int, title string, body string) {
    d.Docs[docID] = Doc{Title: title, Body: body}
}

func (d *DocStore) Get(docID int) Doc {
	return d.Docs[docID]
}

func (d *DocStore) SaveToDisk(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(d.Docs)
}

func (d *DocStore) LoadFromDisk(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&d.Docs)
}
