package link

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

type Link struct {
	Href string `xml:"url>loc"`
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("Link{\n"+
		"  Href: %s\n"+
		"  Text: %s\n"+
		"}", l.Href, l.Text)
}

// getLinks gets all Link types from a .html file
func getLinks(r io.Reader) ([]Link, error) {
	var links []Link

	htmlTokenizer := html.NewTokenizer(r)

	for {
		_ = htmlTokenizer.Next()
		token := htmlTokenizer.Token()

		if token.Type == html.ErrorToken {
			return links, htmlTokenizer.Err()
		}

		// if the token is of type <a>, we need to parse all the content until </a> token is found
		if token.Type == html.StartTagToken && token.Data == "a" {

			var link Link
			link.Href = token.Attr[0].Val

			text, err := parseText(htmlTokenizer)
			if err != nil && err != io.EOF {
				log.Fatalf("error getting text from <a> tag: %s\n", err)
			}
			link.Text = text

			links = append(links, link)
		}
	}
}

// parseText gets all the text between and <a> tag and its ending </a> tag
func parseText(z *html.Tokenizer) (string, error) {
	var text string

	for {
		tokenType := z.Next()

		switch tokenType {
		case html.EndTagToken:
			if z.Token().Data == "a" { // endTag of type </a>, we can return the text
				return text, nil
			}
			// endTag, but not </a> (maybe </i> or </strong>, so we skip it
		case html.ErrorToken: // EOF or something bad happened
			return text, z.Err()
		case html.TextToken: // new text found, so we append this text to the final result
			found := string(z.Text())
			text += cleanString(found)
		}
	}
}

// cleanString cleans the s string from newlines and leading whitespaces (not trailing ones!)
func cleanString(s string) string {
	withoutNewlines := strings.ReplaceAll(s, "\n", "")
	return strings.TrimLeft(withoutNewlines, " ")
}

// ValidLink checks if a link is eligible to be added to a list of links
func (l Link) IsValid(domain string, links []Link) bool {
	if !l.isLinkToSameDomain(domain) || l.linkExists(links) {
		return false
	}
	return true
}

// isLinkToSameDomain checks if a link corresponds to a certain domain
func (l Link) isLinkToSameDomain(domain string) bool {
	s := l.Href
	// this is a redirection to same domain
	if string(s[0]) == "/" {
		return true
	}

	// parse the host
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}

	// the domain does not match
	if !strings.Contains(u.Host, domain) {
		return false
	}

	return true
}

// ParseLink generates a well-formed (with domain) link from any valid link
func (l Link) ParseLink(domain string) Link {
	u, err := url.Parse(l.Href)
	if err != nil {
		log.Fatal(err)
	}

	var tmpLink string

	if len(u.RawQuery) == 0 {
		tmpLink = domain + u.EscapedPath()
	} else {
		tmpLink = domain + u.EscapedPath() + "?" + u.RawQuery
	}

	if strings.Contains(tmpLink, "http") {
		return Link{Href: tmpLink}
	}
	return Link{Href: "https://" + tmpLink}
}

// linkExists returns true if a link exists in a list of links
func (l Link) linkExists(list []Link) bool {
	return slices.Contains(list, l)
}

func GetLinksFromPage(website string) []Link {
	var res *http.Response
	var err error

	if strings.Contains(website, "http") {
		res, err = http.Get(website)
	} else {
		res, err = http.Get("https://" + website)
	}

	if err != nil {
		log.Fatalf("error getting %s: %s\n", website, err)
		return nil
	}

	defer res.Body.Close()

	// here we get the initial links, but we need to visit them and add their links also
	links, err := getLinks(res.Body)
	if err != nil && err != io.EOF {
		log.Fatalf("error getting Link Tokens from website %s: %s\n", website, err)
	}

	var list []Link
	for _, link := range links {
		if link.IsValid(website, list) {
			list = append(list, link.ParseLink(website))
		}
	}

	return list
}
