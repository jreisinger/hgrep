// hgrep is a command-line tool to search URLs for patterns
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/jreisinger/hgrep"
)

var i = flag.Bool("i", false, "perform case insensitive matching")
var m = flag.Bool("m", false, "print only matched parts")
var r = flag.Bool("r", false, "search links within host recursively")
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
	var n int // number of pending sends to worklist to know when we're done

	n++
	go func() {
		worklist <- urls // start with the command-line arguments
	}()

	// tokens is a counting semaphore used to
	// enforce a limit on concurrent requests.
	var tokens = make(chan struct{}, *c) // struct{} has size zero

	seen := make(map[string]bool) // to dedup links
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					tokens <- struct{}{} // acquire a token
					result := hgrep.Grep(link, rx, *r)
					if result.Err != nil {
						log.Printf("greping %s: %v", result.URL, result.Err)
					} else {
						result.Print(*m, headers)
						if *r {
							worklist <- result.Links
						} else {
							worklist <- nil // you need to send something even if it's nothing :-)
						}
					}
					<-tokens // release the token
				}(link)
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
