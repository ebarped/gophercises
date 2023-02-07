package link

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("Link{\n"+
		"  Href: %s\n"+
		"  Text: %s\n"+
		"}", l.Href, l.Text)
}

// GetLinksFromFile gets all Link types from a .html file
func GetLinks(r io.Reader) ([]Link, error) {
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
