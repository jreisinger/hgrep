// Package hgrep (HTTP grep) searches URLs for patterns.
package hgrep

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/jreisinger/hgrep/links"
)

const colorReset = "\033[0m"
const colorBlue = "\033[34m"
const colorRed = "\033[31m"

// Grep searches url for pattern and if found prints the line that matched.
// If recurse is true it also extracts links from the URL and returns them.
// Header enables printing of URLs at which pattern was found. MatchesOnly
// prints only matched parts not whole lines.
// func Grep(url string, pattern *regexp.Regexp, recurse, header, matchesOnly bool) []string {
// 	result := fetchAndSearch(url, pattern, matchesOnly)
// 	result.print(header)

// 	if recurse {
// 		list, err := links.Extract(url, true)
// 		if err != nil {
// 			log.Print(err)
// 		}
// 		return list
// 	}

// 	return nil
// }

// Grep searches url for pattern.
func Grep(url string, pattern *regexp.Regexp, extractLinks bool) Result {
	result := Result{URL: url}

	body, err := fetch(url)
	if err != nil {
		result.Err = err
		return result
	}

	if extractLinks {
		result.Links, result.Err = links.Extract(url, true)
		if result.Err != nil {
			return result
		}
	}

	result.Lines, result.Matches, result.Err = match(body, pattern)
	if result.Err != nil {
		return result
	}

	return result
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Result contains search results.
type Result struct {
	URL     string
	Lines   []string // lines that matched
	Matches []string // only the matched parts
	Links   []string // links from URL with the same host
	Err     error
}

// Print prints the Result lines or matches, with or without URL headers.
func (r Result) Print(matches, headers bool) {
	out := r.Lines
	if matches {
		out = r.Matches
	}
	for _, o := range out {
		if headers {
			fmt.Printf("%s", colorBlue)
			fmt.Printf("%s", r.URL)
			fmt.Printf("%s", colorReset)
			fmt.Printf("%s", ":")
		}
		fmt.Printf("%s\n", o)
	}

}

// match searches input for pattern matches. It returns whole lines with
// highlight matches and matches themselves.
func match(input []byte, pattern *regexp.Regexp) (lines, matches []string, err error) {
	for _, line := range strings.Split(string(input), "\n") {
		// get matches
		ss := pattern.FindAllStringSubmatch(line, -1)
		if len(ss) != 0 {
			for _, s := range ss {
				if pattern.NumSubexp() > 0 {
					matches = append(matches, s[1:]...)
				} else {
					matches = append(matches, s...)
				}
			}
		}

		// get lines with highlights
		is := pattern.FindAllStringIndex(line, -1)
		if is == nil {
			continue
		}
		var highlight string
		var s int
		for _, i := range is {
			highlight += fmt.Sprintf("%s", line[s:i[0]])
			highlight += fmt.Sprintf("%s", colorRed)
			highlight += fmt.Sprintf("%s", line[i[0]:i[1]])
			highlight += fmt.Sprintf("%s", colorReset)
			s = i[1]
		}
		highlight += fmt.Sprintf("%s", line[s:])
		lines = append(lines, highlight)
	}

	return lines, matches, nil
}
