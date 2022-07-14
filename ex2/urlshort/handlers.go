package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectPath, ok := pathsToUrls[r.URL.Path]; ok {
			log.Printf("MapHandler: Path %s matched, redirecting to %s\n", r.URL.Path, redirectPath)
			http.Redirect(w, r, redirectPath, http.StatusPermanentRedirect)
		} else {
			log.Printf("MapHandler: Path %s not matched, falling back to default handler...\n", r.URL.Path)
			fallback.ServeHTTP(w, r) // call original handler
		}

	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type Config struct {
	Path string `yaml:"path",json:"path"`
	URL  string `yaml:"url",json:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var config []Config
	err := yaml.Unmarshal(yml, &config)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, cfg := range config {
		m[cfg.Path] = cfg.URL
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("YAMLHandler: Path %s matched, redirecting to MapHandler to handle it ...\n", r.URL.Path)
		MapHandler(m, fallback).ServeHTTP(w, r)
	}), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//      [
//        {
//          "path": "/test1",
//          "url": "https://google.es"
//        },
//        {
//          "path": "/test2",
//          "url": "https://duckduckgo.com"
//        }
//      ]
// The only errors that can be returned all related to having
// invalid JSON data.
//
func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var config []Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, cfg := range config {
		m[cfg.Path] = cfg.URL
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("JSONHandler: Path %s matched, redirecting to MapHandler to handle it ...\n", r.URL.Path)
		MapHandler(m, fallback).ServeHTTP(w, r)
	}), nil
}

// DBHandler will interact with the underlying database and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the DB, then the
// fallback http.Handler will be called instead.
//
// The only errors that can be returned all related to failing
// to read from the database.
//
func DBHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	m := make(map[string]string)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		b.ForEach(func(k, v []byte) error {
			log.Printf("DBHandler: Registering %s:%s\n", string(k), string(v))
			m[string(k)] = string(v)
			return nil
		})
		return nil
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DBHandler: Path %s matched, redirecting to MapHandler to handle it ...\n", r.URL.Path)
		MapHandler(m, fallback).ServeHTTP(w, r)
	}), nil
}
