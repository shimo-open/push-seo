package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocsFromXml(t *testing.T) {
	cases := []struct {
		xmlUrl   string
		expected []string
	}{
		{
			xmlUrl:   "https://raw.githubusercontent.com/shimo-open/push-seo/master/process/example.xml",
			expected: []string{"https://a.com/path-a/", "https://a.com/path-a/apis"},
		},
	}

	for _, c := range cases {
		// 调用被测试的函数
		urlSet, err := FetchXML(c.xmlUrl)
		assert.NoError(t, err)
		t.Log("urlSet info", urlSet)

		locs := ExtractLocsFromSitemap(urlSet)
		assert.Equal(t, c.expected, locs)
	}
}
