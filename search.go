// Package hgrep (HTTP grep) searches URLs for patterns.
package hgrep

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jreisinger/hgrep/links"
)

const colorReset = "\033[0m"
const colorBlue = "\033[34m"
const colorRed = "\033[31m"

// Search searches url for pattern and if found prints the line that matched.
// If recurse is true it also extracts links from the URL and returns them.
// Header enables printing of URLs at which pattern was found. MatchesOnly
// prints only matched parts not whole lines.
func Search(url string, pattern *regexp.Regexp, recurse, header, matchesOnly bool) []string {
	result := fetchAndSearch(url, pattern, matchesOnly)
	result.print(header)

	if recurse {
		list, err := links.Extract(url, true)
		if err != nil {
			log.Print(err)
		}
		return list
	}

	return nil
}

// result contains the URL that was searched, any lines that matched the
// pattern and possible error.
type result struct {
	url   string
	lines []string
	err   error
}

// fetchAndSearch fetches the url and searches rx in it.
func fetchAndSearch(url string, rx *regexp.Regexp, matchesOnly bool) result {
	result := result{url: url}

	resp, err := http.Get(url)
	if err != nil {
		result.err = err
		return result
	}
	defer resp.Body.Close()

	result.lines, err = match(resp.Body, rx, matchesOnly)
	if err != nil {
		result.err = err
		return result
	}

	return result
}

// print prints the Result in colors.
func (r result) print(headers bool) {
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

// match searches input for pattern matches. It returns lines witch highlighted
// matches or matches only.
func match(input io.Reader, pattern *regexp.Regexp, matchesOnly bool) (lines []string, err error) {
	b, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if matchesOnly {
			matches := pattern.FindAllStringSubmatch(line, -1)
			if len(matches) != 0 {
				for _, m := range matches {
					lines = append(lines, m...)
				}
			}
		} else {
			matches := pattern.FindAllStringIndex(line, -1)
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
