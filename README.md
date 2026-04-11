# 🚀 Zeno — Full-Text Search Engine from Scratch in Go

A full-text search engine built from scratch in Go. Zeno parses Wikipedia XML dumps, builds an inverted index, persists it to disk as binary files, and lets you search across hundreds of thousands of articles in milliseconds — via CLI or HTTP API.

No external dependencies. No databases. Just Go and raw data structures.

---

## 📑 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [How It Works — End to End](#-how-it-works--end-to-end)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [Binary File Format](#-binary-file-format-deep-dive)
- [Project Structure](#-project-structure)
- [Example Output](#-example-output)
- [License](#-license)

---

## ✨ Features

- **Streaming XML parser** — decompresses `.bz2` Wikipedia dumps on-the-fly, no need to extract the full file
- **Tokenization pipeline** — lowercasing, non-alphanumeric stripping, stopword filtering
- **Inverted index** with term-frequency tracking
- **Persistent storage** — index and documents serialized to disk using Go's `encoding/gob`
- **Boolean AND search** — multi-token queries return only docs containing _all_ terms
- **HTTP JSON API** — serve the index over HTTP for integration with frontends/tools
- **Zero external dependencies** — pure Go standard library

---

## 🧠 Architecture

```
Search Engine Pipeline (High-Level Architecture):
```
<img width="1468" height="588" alt="Untitled-2026-02-02-2156" src="https://github.com/user-attachments/assets/6b4c22e8-1142-4b3b-9352-bf32015d488a" />

```
┌──────────────────────────────────────────────────────────────────────┐
│                         INDEXING PHASE                               │
│                                                                      │
│  Wikipedia XML.bz2                                                   │
│       │                                                              │
│       ▼                                                              │
│  ┌──────────┐    ┌────────────┐    ┌───────────┐    ┌─────────────┐  │
│  │  Parser   │───▶│ Tokenizer  │───▶│ Stopword  │───▶│  Inverted   │  │
│  │ (stream)  │    │ (lowercase │    │  Filter   │    │   Index     │  │
│  │           │    │  + clean)  │    │           │    │             │  │
│  └──────────┘    └────────────┘    └───────────┘    └──────┬──────┘  │
│                                                            │         │
│                                                            ▼         │
│                                              ┌─────────────────────┐ │
│                                              │  gob.Encode → Disk  │ │
│                                              │  zeno_bin (index)   │ │
│                                              │  zeno_docs (docs)   │ │
│                                              └─────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│                          SEARCH PHASE                                │
│                                                                      │
│  User Query                                                          │
│       │                                                              │
│       ▼                                                              │
│  ┌────────────┐    ┌───────────┐    ┌──────────────┐                 │
│  │ Tokenizer  │───▶│ Stopword  │───▶│ Index Lookup │                 │
│  │            │    │  Filter   │    │ (AND across  │                 │
│  │            │    │           │    │  all tokens) │                 │
│  └────────────┘    └───────────┘    └──────┬───────┘                 │
│                                            │                         │
│                                            ▼                         │
│                                    ┌──────────────┐                  │
│                                    │   Doc Store   │                  │
│                                    │ (title+body)  │                  │
│                                    └──────┬───────┘                  │
│                                           │                          │
│                                           ▼                          │
│                                   Search Results                     │
│                                   (CLI or JSON)                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 How It Works — End to End

### 1. Parse Wikipedia Dump
The parser (`internal/parser`) opens the `.bz2` compressed XML file and streams `<page>` elements one-by-one using Go's `encoding/xml` decoder. It automatically:
- Decompresses bzip2 on-the-fly (no need to extract the 2GB+ XML)
- Skips **redirect pages** (pages that just point to another page)
- Skips **empty articles** (no revision text)
- Sends each valid `Article{Title, Body}` into a buffered channel (size 100)

### 2. Tokenize
The tokenizer (`internal/tokenizer`) processes each article's full text (`title + body`):
1. **Lowercases** the entire text
2. **Splits** on whitespace into individual tokens
3. **Strips non-alphanumeric** characters via regex `[^a-z0-9]+`
4. **Removes stopwords** (common English words like "the", "is", "and", etc.)

### 3. Build the Inverted Index
The index (`internal/index`) maintains a map:
```
token → { docID → termFrequency }
```
For each document, every surviving token is recorded along with how many times it appears in that document.

### 4. Persist to Disk
Both the inverted index and document store are serialized using Go's `encoding/gob` format:
- `zeno_bin` — the inverted index (~480 MB for Simple English Wikipedia)
- `zeno_docs` — all document titles and bodies (~1.2 GB)

### 5. Search
At query time:
1. Load `zeno_bin` and `zeno_docs` from disk
2. Tokenize the query using the same pipeline
3. For each query token, look up the posting list from the index
4. **Intersect** all posting lists (boolean AND) — only documents containing _every_ query token are returned
5. Look up titles/bodies from the doc store and return results

---

## 🛠 Getting Started

### Prerequisites

| Requirement | Version |
|---|---|
| **Go** | 1.25+ |
| **Disk space** | ~2 GB for the index + docs files |
| **RAM** | ~2–3 GB during indexing |

### Step 1 — Clone the Repository

```bash
git clone https://github.com/CHIRANTAN-001/zeno.git
cd zeno
```

### Step 2 — Download the Wikipedia Dump

Download the **Simple English Wikipedia** XML dump (bz2 compressed, ~330 MB):

```bash
mkdir -p data
```

Download the file from:
> **https://dumps.wikimedia.org/simplewiki/latest/simplewiki-latest-pages-articles.xml.bz2**

Place it inside the `data/` directory so the path looks like:
```
zeno/
  data/
    simplewiki-latest-pages-articles.xml.bz2
```

> [!TIP]
> You can use `wget` or `curl` to download it directly:
> ```bash
> wget -P data/ https://dumps.wikimedia.org/simplewiki/latest/simplewiki-latest-pages-articles.xml.bz2
> ```
> Or with curl:
> ```bash
> curl -L -o data/simplewiki-latest-pages-articles.xml.bz2 https://dumps.wikimedia.org/simplewiki/latest/simplewiki-latest-pages-articles.xml.bz2
> ```

> [!NOTE]
> The compressed `.bz2` file is ~330 MB. You do **not** need to decompress it — Zeno reads it directly.

### Step 3 — Build the Index

```bash
cd cmd
go run main.go index
```

This will:
1. Stream-parse the Wikipedia XML dump
2. Tokenize every article
3. Build the inverted index in memory
4. Save `zeno_bin` (index) and `zeno_docs` (documents) to the project root

> [!IMPORTANT]
> The indexing step takes **15–30 minutes** depending on your machine and produces two binary files in the project root — `zeno_bin` (~480 MB) and `zeno_docs` (~1.2 GB). Make sure you have enough disk space.

### Step 4 — Search!

You can now search using the CLI or the HTTP server.

---

## 🔍 Usage

### CLI Search

```bash
cd cmd
go run main.go search "solar system"
```

Example output:
```
Index loaded in 12.5s
Docs loaded in 8.2s
Search took 1.2ms

Results for "solar system":

[1423] Solar System
[5621] Planet
[8932] Jupiter
...

42 results found
```

### HTTP Server

Start the search server:

```bash
cd cmd
go run main.go serve
```

The server starts on `http://localhost:8080`. Query using:

```bash
curl "http://localhost:8080/search?q=solar+system"
```

Response:
```json
{
  "query": "solar system",
  "count": 42,
  "results": [
    {
      "id": 1423,
      "title": "Solar System",
      "snippet": "The Solar System is the gravitationally bound system of the Sun and the objects..."
    }
  ],
  "took_ms": 1
}
```

---

## 💾 Binary File Format (Deep Dive)

Both binary files use Go's [`encoding/gob`](https://pkg.go.dev/encoding/gob) serialization format. Gob is a self-describing, binary, stream-oriented encoding — think of it as Go's native alternative to Protocol Buffers or MessagePack.

### `zeno_bin` — Inverted Index

**Go type being encoded:**
```go
map[string]map[int]int
```

**Logical structure:**
```
┌─────────────────────────────────────────────────────────────┐
│                     GOB STREAM HEADER                       │
│  (type descriptors for map[string]map[int]int)              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Token: "solar"                                     │    │
│  │  ┌───────────────────────────────────────────────┐   │    │
│  │  │  DocID: 1423  →  TermFrequency: 8             │   │    │
│  │  │  DocID: 5621  →  TermFrequency: 2             │   │    │
│  │  │  DocID: 8932  →  TermFrequency: 1             │   │    │
│  │  │  ...                                          │   │    │
│  │  └───────────────────────────────────────────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Token: "planet"                                    │    │
│  │  ┌───────────────────────────────────────────────┐   │    │
│  │  │  DocID: 5621  →  TermFrequency: 15            │   │    │
│  │  │  DocID: 9201  →  TermFrequency: 3             │   │    │
│  │  │  ...                                          │   │    │
│  │  └───────────────────────────────────────────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ... (one entry per unique token in the entire corpus)      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**What each field means:**

| Field | Type | Description |
|---|---|---|
| Token (key) | `string` | A normalized word (lowercased, cleaned, non-stopword) |
| DocID (inner key) | `int` | Unique document identifier (1-based, assigned sequentially) |
| TermFrequency (inner value) | `int` | Number of times this token appears in this document |

### `zeno_docs` — Document Store

**Go type being encoded:**
```go
map[int]Doc

type Doc struct {
    Title string
    Body  string
}
```

**Logical structure:**
```
┌─────────────────────────────────────────────────────────────┐
│                     GOB STREAM HEADER                       │
│  (type descriptors for map[int]Doc)                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  DocID: 1                                           │    │
│  │  ┌───────────────────────────────────────────────┐   │    │
│  │  │  Title: "April"                               │   │    │
│  │  │  Body:  "April is the fourth month of the..." │   │    │
│  │  └───────────────────────────────────────────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  DocID: 2                                           │    │
│  │  ┌───────────────────────────────────────────────┐   │    │
│  │  │  Title: "August"                              │   │    │
│  │  │  Body:  "August is the eighth month of..."    │   │    │
│  │  └───────────────────────────────────────────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ... (one entry per indexed Wikipedia article)              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Gob Encoding — How It Works at the Byte Level

Go's `gob` format writes data as a self-describing binary stream:

```
┌──────────┬────────────────┬──────────────────────────────────┐
│  LENGTH  │  TYPE ID       │  PAYLOAD                         │
│ (varint) │  (varint, for  │  (field values encoded           │
│          │   first value  │   recursively as gob)            │
│          │   or new type) │                                  │
└──────────┴────────────────┴──────────────────────────────────┘
```

1. **Type descriptors** are sent once at the start — they describe the Go types (`map`, `string`, `int`, `struct`) including field names.
2. **Map entries** are encoded as `[length][key][value]` pairs.
3. **Strings** are encoded as `[byte-length][utf8-bytes]`.
4. **Ints** are encoded as variable-length integers (varints).
5. **Structs** are encoded field-by-field in definition order, with field numbers.

> [!NOTE]
> Because `gob` is a Go-native format, these files can only be read by Go programs. They are not human-readable. If you need cross-language support, the serialization layer can be swapped to JSON, MessagePack, or Protocol Buffers.

---

## 📁 Project Structure

```
zeno/
├── cmd/
│   └── main.go                  # Entry point — CLI & HTTP server
├── data/
│   └── simplewiki-...xml.bz2    # Wikipedia dump (not committed, download separately)
├── internal/
│   ├── index/
│   │   ├── index.go             # Inverted index — build, save, load
│   │   └── docs.go              # Document store — title+body, save, load
│   ├── parser/
│   │   └── wikipedia.go         # Streaming bzip2 XML parser
│   ├── search/
│   │   └── search.go            # Boolean AND search across posting lists
│   ├── stopwords/
│   │   └── stopwords.go         # Stopword filter (common English words)
│   └── tokenizer/
│       └── tokenizer.go         # Tokenization — lowercase, clean, filter
├── zeno_bin                     # Inverted index binary (generated, ~480 MB)
├── zeno_docs                    # Document store binary (generated, ~1.2 GB)
├── go.mod                       # Go module definition
└── README.md                    # This file
```

---

## 🧪 Testing Guide — Step-by-Step

Here's the full flow to test Zeno from scratch on a fresh machine:

```bash
# 1. Clone the repo
git clone https://github.com/CHIRANTAN-001/zeno.git
cd zeno

# 2. Download the Wikipedia dump (~330 MB, takes a few minutes)
mkdir -p data
wget -P data/ https://dumps.wikimedia.org/simplewiki/latest/simplewiki-latest-pages-articles.xml.bz2

# 3. Build the index (15–30 min, generates zeno_bin + zeno_docs)
cd cmd
go run main.go index

# 4a. Test via CLI
go run main.go search "albert einstein"

# 4b. Test via HTTP
go run main.go serve
# In another terminal:
curl "http://localhost:8080/search?q=albert+einstein"
```

**What to expect:**

| Step | Output |
|---|---|
| **Indexing** | Prints `Indexed 10000 articles`, `Indexed 20000 articles`, ... until done |
| **CLI search** | Shows index/docs load time, search time, and matching article titles |
| **HTTP search** | Returns JSON with `query`, `count`, `results` array, and `took_ms` |

---

## 📝 License

This project is for educational purposes — built to understand how real-world search engines like Elasticsearch and Google work under the hood.
