package process

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/ego-component/eos"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/elog"
)

type Processor struct {
	sitemaps []string
	token    string
	s3Cli    *eos.Component
	dryRun   bool
	httpCli  *ehttp.Component
}

// Urlset Struct to unmarshal XML content
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

func NewProcessor(sitemaps []string, token string, dryRun bool, s3Cli *eos.Component, httpCli *ehttp.Component) *Processor {
	return &Processor{
		sitemaps: sitemaps,
		token:    token,
		s3Cli:    s3Cli,
		dryRun:   dryRun,
		httpCli:  httpCli,
	}
}

func (p *Processor) Process() error {
	if p.token == "" || len(p.sitemaps) == 0 {
		return fmt.Errorf("fail: the BAIDU_TOKEN environment variable and URL flag are required")
	}
	for _, v := range p.sitemaps {
		if err := p.ProcessSitemapURL(v, p.token); err != nil {
			return fmt.Errorf("ProcessSitemapURL fail, %w", err)
		}
	}
	return nil
}

const BaiduApi = "http://data.zz.baidu.com/urls?site=%s&token=%s"

func ConstructSubmitURL(sitemap string, token string) string {
	u, err := url.Parse(sitemap)
	if err != nil {
		elog.Warn("Error parsing URL", l.S("sitemap", sitemap), l.E(err))
		return ""
	}
	return fmt.Sprintf(BaiduApi, u.Scheme+"://"+u.Host, token)
}

// ProcessSitemapURL 处理 sitemap 的函数
func (p *Processor) ProcessSitemapURL(sitemap string, token string) error {
	urlset, err := FetchXML(sitemap)
	if err != nil {
		return fmt.Errorf("Error fetching Shimo.im XMLvfrom %s: %w\n", sitemap, err)
	}

	// 提取locs
	locs := ExtractLocsFromSitemap(urlset)

	// 存储并推送
	if err = p.StoreAndPush(sitemap, token, locs); err != nil {
		return fmt.Errorf("error storing locs in EOS: %w\n", err)
	}
	return nil
}

func FetchXML(sitemap string) (Urlset, error) {
	resp, err := http.Get(sitemap)
	if err != nil {
		return Urlset{}, fmt.Errorf("failed to get Urlset%w\n", err)
	}
	defer resp.Body.Close()

	var result Urlset
	err = xml.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// ExtractLocsFromSitemap 提取SitemapUrl中的Locs
func ExtractLocsFromSitemap(urlset Urlset) []string {
	var locs []string
	for _, u := range urlset.Urls {
		locs = append(locs, u.Loc)
	}
	return locs
}

// StoreAndPush 存储并推送 locs 的函数
func (p *Processor) StoreAndPush(sitemap string, token string, locs []string) error {
	// 获取旧的locs
	oldLocs, err := p.GetLocsFromEOS(sitemap)
	if err != nil {
		return fmt.Errorf("error Get old locs from EOS: %w\n", err)
	}

	// 去重与排序新的 locs
	locs = DeduplicateAndSort(locs)

	// 判断新旧locs是否相等
	if !slices.Equal(locs, oldLocs) {
		if err := p.StoreLocsInEOS(sitemap, locs); err != nil {
			return fmt.Errorf("locs = oldlocs: %w\n", err)
		}
		return p.PushToBaiduSEO(sitemap, token, locs)
	}
	return nil
}

// DeduplicateAndSort 去重与排序locs
func DeduplicateAndSort(locs []string) []string {
	// 去重
	locMap := make(map[string]struct{})
	var deduplicatedLocs []string
	for _, loc := range locs {
		if _, found := locMap[loc]; !found {
			locMap[loc] = struct{}{}
			deduplicatedLocs = append(deduplicatedLocs, loc)
		}
	}
	// 排序
	sort.Strings(deduplicatedLocs)
	return deduplicatedLocs
}

// StoreLocsInEOS 将locs保存在EOS对象存储中
func (p *Processor) StoreLocsInEOS(sitemap string, locs []string) error {
	if err := p.SaveLocsInEOS(p.s3Cli, sitemap, locs); err != nil {
		return fmt.Errorf("saving Locs in EOS failed: %w", err)
	}
	return nil
}

// SaveLocsInEOS 将locs保存在EOS对象存储中
func (p *Processor) SaveLocsInEOS(c *eos.Component, sitemap string, locs []string) error {
	locsStr := strings.Join(locs, "\n")
	if err := c.Put(context.Background(), sitemap, strings.NewReader(locsStr), nil); err != nil {
		return fmt.Errorf("EOS put failed: %w\n", err)
	}
	return nil
}

// GetLocsFromEOS 从EOS中取oldLocs
func (p *Processor) GetLocsFromEOS(sitemap string) ([]string, error) {
	oldLocs, err := p.GetOldLocsFromEOS(p.s3Cli, sitemap)
	if err != nil {
		return nil, fmt.Errorf("GetOldLocsFromEOS fail, %w", err)
	}
	return oldLocs, nil
}

// GetOldLocsFromEOS 从EOS中取oldLocs
func (p *Processor) GetOldLocsFromEOS(c *eos.Component, sitemap string) ([]string, error) {
	data, err := c.Get(context.Background(), sitemap)
	if err != nil {
		return nil, fmt.Errorf("oldLocs getting from EOS failed: %v", err)
	}
	oldLocs := strings.Split(data, "")
	return oldLocs, nil
}

// PushToBaiduSEO 推送到百度SEO
func (p *Processor) PushToBaiduSEO(sitemap string, token string, locs []string) error {
	baiduApi := ConstructSubmitURL(sitemap, token)
	locStr := strings.Join(locs, "\n")
	if p.dryRun {
		elog.Info("Push msg with dry run", l.S("locStr", locStr), l.S("api", baiduApi))
		return nil
	}
	res, err := p.httpCli.R().SetBody(strings.NewReader(locStr)).Post(baiduApi)
	if err != nil {
		return fmt.Errorf("error submitting loc %s to Baidu: %w\n", locStr, err)
	}
	elog.Info("Baidu SEO update response", l.S("locStr", locStr), l.A("res", res))

	return nil
}
