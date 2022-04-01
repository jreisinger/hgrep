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
* [x] support `-H` to always print URL headers
* [x] support `-h` to never print URL headers
* [x] supprot `-m` to print only matched parts
* [ ] support `-r` for recursive links search