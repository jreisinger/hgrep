package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// linksExtract makes an HTTP GET request to the specified URL, parses the
// response as HTML, and returns the links in the HTML document. It ignores bad
// URLs and URLs on different host.
func linksExtract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil || !sameHost(url, link) {
					// Ignore bad URLs and
					// URLs on other hosts.
					continue
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func sameHost(URL string, link *url.URL) bool {
	u, err := url.Parse(URL)
	if err != nil {
		return false
	}
	return u.Host == link.Host
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
