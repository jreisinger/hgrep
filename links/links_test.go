package links

import (
	"reflect"
	"testing"
)

func TestSameHost(t *testing.T) {
	testcases := []struct {
		url1, url2 string
		sameHost   bool
	}{
		{"", "", true},
		{" ", " ", true},
		{"a", "b", true}, // empty host
		{"https://example.com", "https://example.com", true},
		{"http://example.com", "https://example.com", true},
		{"https://example.com", "http://example.com", true},
		{"https://example.com", "https://example.net", false},
		{"https://golang.org", "https://go.dev", false},
		{"https://perl.com", "https://go.dev", false},
	}
	for _, tc := range testcases {
		got := sameHost(tc.url1, tc.url2)
		if got != tc.sameHost {
			t.Fatalf("expected same hosts for %s and %s: %t but got: %t", tc.url1, tc.url2, tc.sameHost, got)
		}
	}
}

func TestExtract(t *testing.T) {
	testcases := []struct {
		url          string
		sameHostOnly bool
		links        []string
	}{
		{
			"https://example.com",
			false,
			[]string{"https://www.iana.org/domains/example"},
		},
		{
			"https://example.com",
			true,
			nil,
		},
	}

	for _, tc := range testcases {
		links, err := Extract(tc.url, tc.sameHostOnly)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(links, tc.links) {
			t.Fatalf("extracting from %s got '%v' but expected '%v'", tc.url, links, tc.links)
		}
	}
}
