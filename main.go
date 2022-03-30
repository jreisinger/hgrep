package main

import (
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
		log.Fatal("usage: [pattern] [url ...]")
	}

	pattern := os.Args[1]
	rx, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	urls := os.Args[2:]
	ch := make(chan Result)
	for _, url := range urls {
		go search(url, rx, ch)
	}
	for range urls {
		result := <-ch
		if result.err != nil {
			log.Printf("%v", result.err)
			continue
		}
		print(result)
	}
}

type Result struct {
	url   string
	lines []string
	err   error
}

func print(result Result) {
	colorReset := "\033[0m"
	colorBlue := "\033[34m"

	for _, line := range result.lines {
		fmt.Printf("%s", colorBlue)
		fmt.Printf("%s: ", result.url)
		fmt.Printf("%s", colorReset)
		fmt.Printf("%s\n", line)
	}

}

func search(url string, rx *regexp.Regexp, ch chan Result) {
	result := Result{url: url}

	resp, err := http.Get(url)
	if err != nil {
		result.err = err
		ch <- result
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		result.err = err
		ch <- result
		return
	}

	var lines []string
	for _, line := range strings.Split(string(b), "\n") {
		if rx.MatchString(line) {
			lines = append(lines, line)
		}
	}
	result.lines = lines
	ch <- result
}
