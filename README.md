`hgrep` (HTTP grep) fetches URLs and prints lines matching pattern.

INSTALLATION

```
make install
```

USAGE

```
hgrep '👉|Go' reisinge.net/notes/go/pointers reisinge.net/cv
```

TODO

* [x] support reading HTML from stdin
* [x] highlight matches within lines
* [x] support `-i` for case insensitive matching
* [ ] support `-r` for recursive links search