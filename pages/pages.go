package pages

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

//FetchTitle retrives the title for a webpage located at specified URL.
func FetchTitle(url *url.URL) (string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return "", fmt.Errorf("Unable to fetch web page: %s", err)
	}

	defer resp.Body.Close()
	return getHTMLTitle(resp.Body)
}

func getHTMLTitle(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("Unable to parse HTML: %s", err)
	}

	result, ok := traverse(doc)
	if !ok {
		return result, fmt.Errorf("Unable to find title tag")
	}

	return result, nil
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result, ok := traverse(c); ok {
			return result, ok
		}
	}

	return "", false
}
