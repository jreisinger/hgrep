`hgrep` (HTTP grep) fetches URLs and prints lines matching pattern.

INSTALLATION

```
make install
```

USAGE

```
hgrep -r DevOps https://go.dev
hgrep -m '\w+\[\.\][a-z]{1,3}' https://blog.google/threat-analysis-group/exposing-initial-access-broker-ties-conti
```

TODO

* [x] support reading URLs from stdin
* [x] highlight matches within lines
* [x] support `-i` to perform case insensitive matching
* [x] support `-m` to print only matched parts
* [x] support `-r` to search links recursively within the host
