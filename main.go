package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [pattern] [url ...]\n", os.Args[0])
		os.Exit(1)
	}

	pattern := os.Args[1]
	rx, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	urls := os.Args[2:]
	ch := make(chan Result)

	if len(urls) == 0 {
		input := bufio.NewScanner(os.Stdin)
		var nLines int
		for input.Scan() {
			nLines++
			go fetchAndMatch(input.Text(), rx, ch)
		}
		for i := 0; i < nLines; i++ {
			result := <-ch
			if result.err != nil {
				log.Printf("%v", result.err)
				continue
			}
			print(result.url, result.lines)
		}
		if err := input.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, url := range urls {
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
}

type Result struct {
	url   string
	lines []string
	err   error
}

func print(url string, lines []string) {
	colorReset := "\033[0m"
	colorBlue := "\033[34m"

	for _, line := range lines {
		fmt.Printf("%s", colorBlue)
		fmt.Printf("%s: ", url)
		fmt.Printf("%s", colorReset)
		fmt.Printf("%s\n", line)
	}

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

func match(input io.ReadCloser, rx *regexp.Regexp) (lines []string, err error) {
	b, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if rx.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
