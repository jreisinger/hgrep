`hgrep` (HTTP grep) fetches URLs and prints lines matching pattern.

INSTALLATION

```
make install
```

USAGE

```
hgrep '[Gg]o' https://reisinge.net/about https://reisinge.net/cv

echo -e "https://reisinge.net/about\nhttps://reisinge.net/cv" | hgrep '[Gg]o'

echo -e "https://reisinge.net/about\nhttps://reisinge.net/cv" > /tmp/urls.txt
cat /tmp/urls.txt | hgrep '[Gg]o'
```

TODO

* [x] support reading HTML from stdin
* [ ] highlight matches within lines
* [ ] support `-r` for recursive search