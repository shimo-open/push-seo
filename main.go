package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func main() {
	// 解析 shimo.im 的 XML
	shimoURL := "https://open.shimo.im/sitemap.xml"
	shimoXML, err := fetchXML(shimoURL)
	if err != nil {
		fmt.Println("Error fetching Shimo.im XML:", err)
		return
	}
	shimoContent := parseAndGetURLs(shimoXML)

	// 提交接口请求 URL
	submitURL := "http://data.zz.baidu.com/urls?site=https://open.shimo.im&token=wX81OpgvaRbEysIT"

	// 向百度提交 URL
	resp, err := http.Post(submitURL, "text/plain", bytes.NewBuffer([]byte(shimoContent)))
	if err != nil {
		fmt.Println("Error submitting URLs to Baidu:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Baidu SEO update response:", string(body))
}

func fetchXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func parseAndGetURLs(xmlData []byte) string {
	var urlSet URLSet
	err := xml.Unmarshal(xmlData, &urlSet)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return ""
	}

	return "https://open.shimo.im/docs/SDK-3.10/concepts/base"

	var urlsBuffer bytes.Buffer
	for _, url := range urlSet.Urls {
		urlsBuffer.WriteString(url.Loc + "\n")
	}
	return urlsBuffer.String()
}
