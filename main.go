package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Struct to unmarshal XML content
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

const MaxURLsPerSubmit = 1

func main() {
	token := os.Getenv("BAIDU_TOKEN")
	var sitemapURLs flagArray
	flag.Var(&sitemapURLs, "sitemapURL", "URL for the sitemap from shimo.im")

	flag.Parse()

	if token == "" || len(sitemapURLs) == 0 {
		fmt.Println("Error: the BAIDU_TOKEN environment variable and sitemapURL flag are required")
		os.Exit(1)
	}

	for _, sitemapURL := range sitemapURLs {
		submitURL := constructSubmitURL(sitemapURL, token)
		shimoXML, err := fetchXML(sitemapURL)
		if err != nil {
			fmt.Printf("Error fetching Shimo.im XML from %s: %v\n", sitemapURL, err)
			continue
		}

		urls, err := parseAndGetURLs(shimoXML)
		if err != nil {
			fmt.Printf("Error parsing Shimo.im XML from %s: %v\n", sitemapURL, err)
			continue
		}

		if len(urls) < 1 {
			fmt.Printf("No URLs found in the sitemap: %s\n", sitemapURL)
			continue
		}

		submitContent := urls[2] // Sending only the third URL from each sitemap

		// Submit to Baidu
		resp, err := http.Post(submitURL, "text/plain", bytes.NewBuffer([]byte(submitContent.Loc)))
		if err != nil {
			fmt.Printf("Error submitting URL to Baidu: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Baidu SEO update response for %s: %s\n", sitemapURL, string(body))
	}
}

type flagArray []string

func (i *flagArray) String() string {
	return ""
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func fetchXML(url string) (Urlset, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Urlset{}, err
	}
	defer resp.Body.Close()

	var result Urlset
	err = xml.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func parseAndGetURLs(sitemap Urlset) ([]Url, error) {
	return sitemap.Urls, nil
}

func constructSubmitURL(sitemapURL, token string) string {
	u, err := url.Parse(sitemapURL)
	if err != nil {
		fmt.Println("Error parsing sitemapURL:", err)
		return ""
	}
	return fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", u.Host, token)
}
