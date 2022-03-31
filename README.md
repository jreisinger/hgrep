`hgrep` (HTTP grep) fetches URLs and prints lines matching pattern.

INSTALLATION

```
make install
```

USAGE

```
hgrep '[Gg][Oo]' reisinge.net/about reisinge.net/cv
echo -e "reisinge.net/about\nreisinge.net/cv" | hgrep -i 'go'
```

TODO

* [x] support reading HTML from stdin
* [x] highlight matches within lines
* [x] support `-i` for case insensitive matching
* [ ] support `-r` for recursive links search