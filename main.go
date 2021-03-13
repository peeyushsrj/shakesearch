package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	// "index/suffixarray"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
)

var index, _ = bleve.Open("completeworks2.bleve")

func main() {
	// searcher := Searcher{}
	// err := searcher.Load("completeworks.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// type Searcher struct {
// 	CompleteWorks string
// 	SuffixArray   *suffixarray.Index
// }

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// return func(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["q"]
	fmt.Println(query)
	if !ok || len(query[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing search query in URL params"))
		return
	}
	searchquery := bleve.NewMatchQuery(query[0])
	search := bleve.NewSearchRequest(searchquery)
	search.Fields = []string{"*"}
	searchResults, err := index.Search(search)

	if err != nil {
		fmt.Println("!Bleeve search panic!")
		fmt.Println(err)
		return
	}
	// fmt.Println(searchResults)

	results := make([]interface{}, 0, len(searchResults.Hits))
	for _, el := range searchResults.Hits {
		results = append(results, el.Fields[""].(string))
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err = enc.Encode(results)
	if err != nil {
		fmt.Println("!Encoder panic!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("encoding failure"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
	// }
}

// func (s *Searcher) Load(filename string) error {
// 	dat, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return fmt.Errorf("Load: %w", err)
// 	}
// 	s.CompleteWorks = string(dat)
// 	s.SuffixArray = suffixarray.New(dat)
// 	return nil
// }

// func (s *Searcher) Search(query string) []string {
// 	idxs := s.SuffixArray.Lookup([]byte(query), -1)
// 	results := []string{}
// 	for _, idx := range idxs {
// 		results = append(results, s.CompleteWorks[idx-250:idx+250])
// 	}
// 	return results
// }
