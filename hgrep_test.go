package hgrep

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	testcases := []struct {
		input []byte
		rx    *regexp.Regexp
		lines []string
	}{
		{
			[]byte(""),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset},
		},
		{
			[]byte(" "),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset + " " + colorRed + colorReset},
		},
		{
			[]byte("a"),
			regexp.MustCompile(``),
			[]string{colorRed + colorReset + "a" + colorRed + colorReset},
		},
		{
			[]byte("a"),
			regexp.MustCompile(`b`),
			nil,
		},
		{
			[]byte("line1\nline2"),
			regexp.MustCompile(`3`),
			nil,
		},
		{
			[]byte("line1\nline2"),
			regexp.MustCompile(`1`),
			[]string{
				"line" + colorRed + "1" + colorReset,
			},
		},
		{
			[]byte("line1\nline2"),
			regexp.MustCompile(`[12]`),
			[]string{
				"line" + colorRed + "1" + colorReset,
				"line" + colorRed + "2" + colorReset,
			},
		},
		{
			[]byte("word1word2"),
			regexp.MustCompile(`1`),
			[]string{
				"word" + colorRed + "1" + colorReset + "word2",
			},
		},
		{
			[]byte("word1word2"),
			regexp.MustCompile(`[12]`),
			[]string{
				"word" + colorRed + "1" + colorReset + "word" + colorRed + "2" + colorReset,
			},
		},
	}
	for _, tc := range testcases {
		lines, _, err := match(tc.input, tc.rx)
		assert.NoError(t, err)
		assert.Equal(t, tc.lines, lines)
	}
}
