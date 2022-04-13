`hgrep` (HTTP grep) searches URLs for patterns.

INSTALLATION

```
make install
```

USAGE

```
hgrep DevOps https://go.dev https://golang.org
hgrep -m '\w+\[\.\][a-z]{1,63}' https://blog.google/threat-analysis-group/rss
hgrep -r DevOps https://go.dev
```

TODO

* [x] support reading URLs from stdin
* [x] highlight matches within lines
* [x] support `-i` to perform case insensitive matching
* [x] support `-m` to print only matched parts
* [x] support `-r` to search links recursively within the host
* [x] support `-c` to limit the number of concurrent searches
