`hgrep` (HTTP grep) fetches URLs and prints lines matching pattern.

INSTALLATION

```
go install
```

USAGE

```
hgrep 'Go' https://reisinge.net/about https://reisinge.net/cv
echo -e "https://reisinge.net/about\nhttps://reisinge.net/cv" | hgrep 'Go'
```

TODO

* [x] support reading HTML from stdin
* [ ] support `-r` for recursive search
* [ ] highlight matches within lines