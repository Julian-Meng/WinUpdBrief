package updates

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	reKB = regexp.MustCompile(`\bKB\d{7}\b`)
)

// LatestEntry
type LatestEntry struct {
	HistoryURL string // update history 入口页
	Title      string // 链接文本（通常含日期 + KB + OS Build）
	KB         string // KB5074109 之类
	URL        string // KB 详情页链接
}

// FetchLatestForBuild 根据 Windows build 自动选择 update history 页面并抓取最新条目
func (f *Fetcher) FetchLatestForBuild(ctx context.Context, build int) (*LatestEntry, error) {
	historyURL := historyURLByBuild(build)
	if historyURL == "" {
		return nil, fmt.Errorf("unsupported windows build: %d", build)
	}
	return f.FetchLatestFromHistory(ctx, historyURL)
}

// 内部：build → history URL 映射
func historyURLByBuild(build int) string {
	// Win11 25H2: 26200+
	if build >= 26200 {
		return "https://support.microsoft.com/en-us/topic/windows-11-version-25h2-update-history-99c7f493-df2a-4832-bd2d-6706baa0dec0"
	}
	// Win11 24H2: 26100+
	if build >= 26100 {
		return "https://support.microsoft.com/en-us/topic/windows-11-version-24h2-update-history-0929c747-1815-4543-8461-0160d16f15e5"
	}
	// Win11 23H2 / 22H2
	if build >= 22000 {
		return "https://support.microsoft.com/en-us/topic/windows-11-update-history-93345c32-4ae1-6d1c-f885-6c0b718ad0e9"
	}
	// Win10
	if build >= 10240 {
		return "https://support.microsoft.com/en-us/topic/windows-10-update-history-8127c2c6-6edf-4fdf-8b9f-0f7be1ef3562"
	}
	return ""
}

// 真正的页面抓取 + 解析
func (f *Fetcher) FetchLatestFromHistory(ctx context.Context, historyURL string) (*LatestEntry, error) {
	html, err := f.GetHTML(ctx, historyURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	base, _ := url.Parse(historyURL)

	var bestHref, bestText, bestKB string

	// 策略：找第一个包含 KBxxxxxxx 的链接
	doc.Find("a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		text := strings.TrimSpace(a.Text())
		if text == "" {
			return true
		}

		kb := reKB.FindString(text)
		if kb == "" {
			return true
		}

		href, ok := a.Attr("href")
		if !ok || href == "" {
			return true
		}
		if !strings.Contains(href, "support.microsoft.com") && !strings.Contains(href, "/topic/") {
			return true
		}

		u, err := url.Parse(href)
		if err != nil {
			return true
		}

		bestHref = base.ResolveReference(u).String()
		bestText = text
		bestKB = kb
		return false
	})

	if bestHref == "" {
		return nil, fmt.Errorf("no KB entry found on history page")
	}

	return &LatestEntry{
		HistoryURL: historyURL,
		Title:      bestText,
		KB:         bestKB,
		URL:        bestHref,
	}, nil
}
