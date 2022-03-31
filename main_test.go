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
			strings.NewReader(""),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset},
		},
		{
			strings.NewReader(" "),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset + " " + colorRed + colorReset},
		},
		{
			strings.NewReader("a"),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset + "a" + colorRed + colorReset},
		},
		{
			strings.NewReader("a"),
			regexp.MustCompile(`b`),
			nil,
		},
		{
			strings.NewReader("line1\nline2"),
			regexp.MustCompile(`3`),
			nil,
		},
		{
			strings.NewReader("line1\nline2"),
			regexp.MustCompile(`1`),
			[]string{
				"line" + colorRed + "1" + colorReset,
			},
		},
		{
			strings.NewReader("line1\nline2"),
			regexp.MustCompile(`[12]`),
			[]string{
				"line" + colorRed + "1" + colorReset,
				"line" + colorRed + "2" + colorReset,
			},
		},
		{
			strings.NewReader("word1word2"),
			regexp.MustCompile(`1`),
			[]string{
				"word" + colorRed + "1" + colorReset + "word2",
			},
		},
		{
			strings.NewReader("word1word2"),
			regexp.MustCompile(`[12]`),
			[]string{
				"word" + colorRed + "1" + colorReset + "word" + colorRed + "2" + colorReset,
			},
		},
	}
	for _, tc := range testcases {
		lines, err := match(tc.input, tc.rx)
		assert.NoError(t, err)
		assert.Equal(t, tc.matches, lines)
	}
}
