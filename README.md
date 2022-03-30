`hgrep` is HTML (or HTTP) grep.

Usage examples:

```
> hgrep "Go" https://reisinge.net/about https://reisinge.net/cv
> curl -s https://reisinge.net/about https://reisinge.net/cv 2>&1 | hgrep "Go"
```

TODO

* [x] support reading HTML from stdin
* [ ] support `-r` for recursive search