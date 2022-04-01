package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const colorReset = "\033[0m"
const colorBlue = "\033[34m"
const colorRed = "\033[31m"

func main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	rx, urls, err := parseCLIargs()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan Result)
	for _, url := range urls {
		url := addScheme(url)
		go fetchAndMatch(url, rx, ch)
	}
	for range urls {
		result := <-ch
		if result.err != nil {
			log.Printf("%v", result.err)
			continue
		}
		print(result.url, result.lines)
	}
}

func parseCLIargs() (rx *regexp.Regexp, urls []string, err error) {
	i := flag.Bool("i", false, "perform case insensitive matching")

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

func addScheme(url string) string {
	if url != "" && !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	return url
}

func print(url string, lines []string) {
	for _, line := range lines {
		fmt.Printf("%s", colorBlue)
		fmt.Printf("%s", url)
		fmt.Printf("%s", colorReset)
		fmt.Printf("%s", ":")
		fmt.Printf("%s\n", line)
	}

}

type Result struct {
	url   string
	lines []string
	err   error
}

func fetchAndMatch(url string, rx *regexp.Regexp, ch chan Result) {
	result := Result{url: url}

	resp, err := http.Get(url)
	if err != nil {
		result.err = err
		ch <- result
		return
	}
	defer resp.Body.Close()

	result.lines, err = match(resp.Body, rx)
	if err != nil {
		result.err = err
		ch <- result
		return
	}

	ch <- result
}

func match(input io.Reader, rx *regexp.Regexp) (lines []string, err error) {
	b, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(b), "\n") {
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

	return lines, nil
}
