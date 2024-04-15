package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var Xml = "https://raw.githubusercontent.com/shimo-open/push-seo/v0.0.2/config/example.yml"

func TestGetLocsFromXml(t *testing.T) {
	// 调用被测试的函数
	urlSet, err := FetchXML(Xml)
	assert.NoError(t, err)
	t.Log("urlSet info", urlSet)

	locs := ExtractLocsFromSitemap(urlSet)
	expected := []string{"https://a.com/path-a/", "https://a.com/path-a/apis"}
	assert.Equal(t, expected, locs)
}
