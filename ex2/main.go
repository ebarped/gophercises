package main

import (
	"ex2/urlshort"
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HelloWorldHandler: Path %s matched\n", r.URL.Path)
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	mux := http.NewServeMux()
	log.Printf("helloWorldHandler: registering %q handler\n", "/")
	mux.HandleFunc("/", helloWorldHandler)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, http.HandlerFunc(helloWorldHandler))
	for url := range pathsToUrls {
		log.Printf("mapHandler: registering %q handler\n", url)
		mux.HandleFunc(url, mapHandler)
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), http.HandlerFunc(helloWorldHandler))
	if err != nil {
		panic(err)
	}

	urls := []string{"/urlshort", "/urlshort-final"}
	for _, url := range urls {
		log.Printf("yamlHandler: registering %q handler\n", url)
		mux.HandleFunc(url, yamlHandler)
	}

	json := `
[
  {
    "path": "/test1",
    "url": "https://google.es"
  },
  {
    "path": "/test2",
    "url": "https://duckduckgo.com"
  }
]
`

	jsonHandler, err := urlshort.JSONHandler([]byte(json), http.HandlerFunc(helloWorldHandler))
	if err != nil {
		panic(err)
	}

	urlsJson := []string{"/test1", "/test2"}
	for _, url := range urlsJson {
		log.Printf("jsonHandler: registering %q handler\n", url)
		mux.HandleFunc(url, jsonHandler)
	}

	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("test"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// insert some path:redirections inside DB on bucket test
	urlsDB := map[string]string{
		"/db1": "https://reddit.com",
		"/db2": "https://muylinux.com",
	}
	for k, v := range urlsDB {
		log.Printf("dbHandler: registering %q handler to redirect to %s\n", k, v)
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("test"))
			err := b.Put([]byte(k), []byte(v))
			return err
		})
	}

	dbHandler, err := urlshort.DBHandler(db, http.HandlerFunc(helloWorldHandler))
	if err != nil {
		panic(err)
	}

	for k, _ := range urlsDB {
		mux.HandleFunc(k, dbHandler)
	}

	log.Println("Starting the server on :8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
