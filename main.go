package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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
