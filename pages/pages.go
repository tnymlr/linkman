package pages

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func FetchTitle(url *url.URL) (string, error) {
	if resp, err := http.Get(url.String()); err == nil {
		defer resp.Body.Close()
		return getHtmlTitle(resp.Body)
	} else {
		return "", fmt.Errorf("Unable to fetch web page: %s", err)
	}
}

func getHtmlTitle(r io.Reader) (string, error) {
	if doc, err := html.Parse(r); err == nil {
		if result, ok := traverse(doc); ok {
			return result, nil
		} else {
			return result, fmt.Errorf("Unable to find title tag")
		}
	} else {
		return "", fmt.Errorf("Unable to parse HTML: %s", err)
	}
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
