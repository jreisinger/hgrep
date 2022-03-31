package main

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	testcases := []struct {
		input   io.Reader
		rx      *regexp.Regexp
		matches []string
	}{
		{
			strings.NewReader("line1\nline2\nline3"),
			regexp.MustCompile(``),
			[]string{"line1", "line2", "line3"},
		},
		{
			strings.NewReader("line1\nline2\nline3"),
			regexp.MustCompile(`line`),
			[]string{"line1", "line2", "line3"},
		},
		{
			strings.NewReader("line1\nline2\nline3"),
			regexp.MustCompile(`1`),
			[]string{"line1"},
		},
		{
			strings.NewReader("line1\nline2\nline3"),
			regexp.MustCompile(`[12]`),
			[]string{"line1", "line2"},
		},
	}
	for _, tc := range testcases {
		lines, err := match(tc.input, tc.rx)
		assert.NoError(t, err)
		assert.Equal(t, tc.matches, lines)
	}
}
