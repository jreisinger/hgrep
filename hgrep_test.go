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

var testHtml = `<ul>
<li><strong>Project Repository:</strong> <a href="https://github.com/knative">https://github.com/knative</a></li>
<li><strong>Contributor Guide:</strong> 
<a href="https://knative.dev/docs/community/contributing/" target="_blank">Knative contributor guide</a></li>
<li><strong>Chat:</strong> 
<a href="https://knative.slack.com/" target="_blank">Knative slack</a></li>
<li><strong>License:</strong> 
<a href="https://choosealicense.com/licenses/apache-2.0/" target="_blank">Apache 2.0</a></li>
<li><strong>Legal Requirements:</strong> 
<a href="https://cla.developers.google.com" target="_blank">Google CLA</a></li>
</ul>`

var testHtml2 = `<ul>
<li><strong>Project Repository:</strong> <a href="https://github.com/kubernetes/kubernetes">https://github.com/kubernetes/kubernetes</a></li>
<li><strong>Contributor Guide:</strong> 
<a href="https://github.com/kubernetes/community/tree/master/contributors/guide" target="_blank">kubernetes/community/contributors/guide</a></li>
<li><strong>Chat:</strong> Slack: 
<a href="https://slack.k8s.io/" target="_blank">slack.k8s.io</a></li>
<li><strong>Developer List/Forum:</strong> 
<a href="https://groups.google.com/forum/#!forum/kubernetes-dev" target="_blank">Kubernetes-dev Mailing List</a></li>
<li><strong>License:</strong> 
<a href="https://choosealicense.com/licenses/apache-2.0/" target="_blank">Apache 2.0</a></li>
<li><strong>Legal Requirements:</strong> 
<a href="https://github.com/cncf/cla" target="_blank">CNCF CLA</a></li>
</ul>`

func TestMatch_matches(t *testing.T) {
	tests := []struct {
		input []byte
		rx    *regexp.Regexp
		want  []string
	}{
		{
			[]byte(testHtml),
			regexp.MustCompile(`"https://github.com/[^"#]+"`),
			[]string{"\"https://github.com/knative\""},
		},
		{
			[]byte(testHtml),
			regexp.MustCompile(`"(https://github.com/[^"#]+)"`),
			[]string{"https://github.com/knative"},
		},
		{
			[]byte(testHtml),
			regexp.MustCompile(`"(https://github.com/([^"#]+))"`),
			[]string{"https://github.com/knative", "knative"},
		},
		{
			[]byte(testHtml2),
			regexp.MustCompile(`"https://github.com/[^"#]+"`),
			[]string{"\"https://github.com/kubernetes/kubernetes\"", "\"https://github.com/kubernetes/community/tree/master/contributors/guide\"", "\"https://github.com/cncf/cla\""},
		},
		{
			[]byte(testHtml2),
			regexp.MustCompile(`"(https://github.com/([^"#]+))"`),
			[]string{"https://github.com/kubernetes/kubernetes", "kubernetes/kubernetes", "https://github.com/kubernetes/community/tree/master/contributors/guide", "kubernetes/community/tree/master/contributors/guide", "https://github.com/cncf/cla", "cncf/cla"},
		},
	}
	for _, test := range tests {
		_, matches, err := match(test.input, test.rx)
		assert.NoError(t, err)
		assert.Equal(t, test.want, matches)
	}
}
