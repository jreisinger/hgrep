package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const colorReset = "\033[0m"
const colorBlue = "\033[34m"
const colorRed = "\033[31m"

var i = flag.Bool("i", false, "perform case insensitive matching")
var m = flag.Bool("m", false, "print only matched parts")
var r = flag.Bool("r", false, "search links recursively within the host")
var c = flag.Int("c", 5, "number of concurrent searches")

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

func main() {
	rx, urls, err := parseCLIargs()
	if err != nil {
		log.Fatal(err)
	}

	var headers bool
	if len(urls) > 1 || *r {
		headers = true
	}

	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- urls }()

	// tokens is a counting semaphore used to
	// enforce a limit on concurrent requests.
	var tokens = make(chan struct{}, *c) // struct{} has size zero

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		for _, url := range <-worklist {
			if !seen[url] {
				seen[url] = true
				n++
				go func(link string) {
					tokens <- struct{}{} // acquire a token
					worklist <- searchAndPrint(link, rx, *r, headers)
					<-tokens // release the token
				}(url)
			}
		}
	}
}

// searchAndPrint searches url for pattern and prints it. If recurse is true it
// searches for links within the URL and returns them. Headers enables printing
// of URLs at which pattern was found.
func searchAndPrint(url string, pattern *regexp.Regexp, recurse, headers bool) []string {
	result := fetchAndSearch(url, pattern)
	result.print(headers)

	if recurse {
		list, err := linksExtract(url)
		if err != nil {
			log.Print(err)
		}
		return list
	}

	return nil
}

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

func parseCLIargs() (rx *regexp.Regexp, urls []string, err error) {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		e := fmt.Errorf("usage: hgrep [flags] [pattern] [url ...]")
		return nil, nil, e
	}

	pattern := args[0]
	if *i {
		pattern = `(?i)` + pattern
	}
	rx, err = regexp.Compile(pattern)
	if err != nil {
		return nil, nil, err
	}

	urls = args[1:]
	if len(urls) == 0 { // get URLs from stdin
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			urls = append(urls, input.Text())
		}
		if err := input.Err(); err != nil {
			return nil, nil, err
		}
	}

	return rx, urls, err
}

type Result struct {
	url   string
	lines []string
	err   error
}

func fetchAndSearch(url string, rx *regexp.Regexp) Result {
	result := Result{url: url}

	resp, err := http.Get(url)
	if err != nil {
		result.err = err
		return result
	}
	defer resp.Body.Close()

	result.lines, err = match(resp.Body, rx)
	if err != nil {
		result.err = err
		return result
	}

	return result
}

func (r Result) print(headers bool) {
	for _, line := range r.lines {
		if headers {
			fmt.Printf("%s", colorBlue)
			fmt.Printf("%s", r.url)
			fmt.Printf("%s", colorReset)
			fmt.Printf("%s", ":")
		}
		fmt.Printf("%s\n", line)
	}

}

func match(input io.Reader, rx *regexp.Regexp) (lines []string, err error) {
	b, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if *m {
			matches := rx.FindAllStringSubmatch(line, -1)
			if len(matches) != 0 {
				for _, m := range matches {
					lines = append(lines, m...)
				}
			}
		} else {
			matches := rx.FindAllStringIndex(line, -1)
			if matches == nil {
				continue
			}
			var highlight string
			var s int
			for _, m := range matches {
				highlight += fmt.Sprintf("%s", line[s:m[0]])
				highlight += fmt.Sprintf("%s", colorRed)
				highlight += fmt.Sprintf("%s", line[m[0]:m[1]])
				highlight += fmt.Sprintf("%s", colorReset)
				s = m[1]
			}
			highlight += fmt.Sprintf("%s", line[s:])
			lines = append(lines, highlight)
		}
	}

	return lines, nil
}
