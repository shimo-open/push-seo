package process

import (
	"fmt"
	"reflect"
	"testing"
)

// 此Xml有时效性，测试时注意更换link
var Xml = "https://raw.githubusercontent.com/shimo-open/push-seo/08158f621474c160994f880d48b433a977d3d5d3/resources/example.xml"

func TestGetLocsFromXml(t *testing.T) {
	// 调用被测试的函数
	urlSet, err := FetchXML(Xml)
	fmt.Printf("urlSet------%v\n", urlSet)
	if err != nil {
		t.Errorf("Error fetching XML: %v", err)
		return
	}

	locs := ExtractLocsFromSitemap(urlSet)
	fmt.Printf("locs--------%v\n", locs)
	// 预期结果
	expected := []string{"https://aa.bb.com/", "https://aa.bb.com/co/"}

	// 比较结果和期望
	if !reflect.DeepEqual(locs, expected) {
		t.Errorf("Mismatch in sitemap URLs. Expected: %v, Actual: %v", expected, locs)
	}
}
