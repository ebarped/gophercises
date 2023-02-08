package main

import (
	"flag"
	"fmt"

	"ex5/pkg/sitemap"
)

func main() {
	website := flag.String("website", "https://google.es", "website to analyze")
	depth := flag.Int("depth", 1, "depth of search")
	flag.Parse()

	fmt.Printf("Website to analyze: %s, Depth: %d\n", *website, *depth)

	links := sitemap.Bfs(*website, *depth)
	sm := sitemap.NewSitemap(*website, links)
	fmt.Println(sm.ToXML())
}
