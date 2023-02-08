package sitemap

import (
	"encoding/xml"
	"fmt"
	"log"

	"ex5/pkg/link"
)

type loc struct {
	Link string `xml:"loc"`
}

func NewLoc(s string) loc {
	return loc{
		Link: s,
	}
}

// LinksToLoc transforms a list of links to a list of locs
func LinksToLoc(links []link.Link) []loc {
	var locs []loc

	for _, link := range links {
		locs = append(locs, NewLoc(link.Href))
	}

	return locs
}

func (l loc) String() string {
	return l.Link
}

type Sitemap struct {
	XMLName xml.Name
	Locs    []loc `xml:"url"`
}

func NewSitemap(domain string, links []link.Link) *Sitemap {
	// transform all links to loc
	locs := LinksToLoc(links)

	return &Sitemap{
		XMLName: xml.Name{Local: "urlset", Space: `http://www.sitemaps.org/schemas/sitemap/0.9`},
		Locs:    locs,
	}
}

// Breadth-first_search
func Bfs(website string, maxDepth int) []link.Link {
	seen := make(map[string]struct{}) // map of already visited links
	var q map[string]struct{}         // queue of links to visit
	nq := map[string]struct{}{        // next-qeue, links generated from the current queue. we will swap this nq to q
		website: {}, // initially, q is empty, nq has only the root website
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		if len(q) == 0 { // if the current level to analyze is empty, we can finish
			break
		}

		for url := range q {
			if _, ok := seen[url]; ok { // if the url has been visited, skip it
				continue
			}
			seen[url] = struct{}{} // add the url to the seen ones

			for _, link := range link.GetLinksFromPage(url) { // get all the links from that url and add them to the next-queue
				nq[link.Href] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(seen)) // add all the seen urls to the result

	for url := range seen {
		result = append(result, url)
	}

	// convert the result to Link struct
	var links []link.Link
	for _, l := range result {
		links = append(links, link.Link{Href: l})
	}

	return links
}

func (s Sitemap) String() string {
	var result string

	for i, link := range s.Locs {
		result += fmt.Sprintf("%d - %s\n", i, link)
	}

	return result
}

func (s *Sitemap) ToXML() string {
	var result string

	result += xml.Header

	out, err := xml.MarshalIndent(s, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	result += string(out) + "\n"

	return result
}
