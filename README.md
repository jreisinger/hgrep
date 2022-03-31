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
* [x] highlight matches within lines
* [x] support `-i` for case insensitive matching
* [ ] support `-r` for recursive links search