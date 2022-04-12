package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jreisinger/hgrep/links"
)

// searchAndPrint searches url for pattern and prints it. If recurse is true it
// searches for links within the URL and returns them. Headers enables printing
// of URLs at which pattern was found.
func searchAndPrint(url string, pattern *regexp.Regexp, recurse, headers bool) []string {
	result := fetchAndSearch(url, pattern)
	result.print(headers)

	if recurse {
		list, err := links.Extract(url, true)
		if err != nil {
			log.Print(err)
		}
		return list
	}

	return nil
}

// Results contains the URL that was searched, any lines that matched the
// pattern and possible error.
type Result struct {
	url   string
	lines []string
	err   error
}

// fetchAndSearch fetches the url and searches rx in it.
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

// print prints the Result in colors.
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

// match searches for rx matches in input.
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
