package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var reKB = regexp.MustCompile(`\bKB\d{7}\b`)

var historyPages = []historyPage{
	{minBuild: 26200, version: "25H2", url: "https://support.microsoft.com/en-us/topic/windows-11-version-25h2-update-history-99c7f493-df2a-4832-bd2d-6706baa0dec0"},
	{minBuild: 26100, version: "24H2", url: "https://support.microsoft.com/en-us/topic/windows-11-version-24h2-update-history-0929c747-1815-4543-8461-0160d16f15e5"},
	{minBuild: 22631, version: "23H2", url: "https://support.microsoft.com/en-us/topic/windows-11-version-23h2-update-history-59875222-b990-4bd9-932f-91a5954de434"},
	{minBuild: 22621, version: "22H2", url: "https://support.microsoft.com/en-us/topic/windows-11-version-22h2-update-history-ec4229c3-9c5f-4e75-9d6d-9025ab70fcce"},
	{minBuild: 22000, version: "21H2", url: "https://support.microsoft.com/en-us/topic/windows-11-version-21h2-update-history-a19cd327-b57f-44b9-84e0-26ced7109ba9"},
	{minBuild: 10240, version: "10", url: "https://support.microsoft.com/en-us/topic/windows-10-update-history-8127c2c6-6edf-4fdf-8b9f-0f7be1ef3562"},
}

type historyPage struct {
	minBuild int
	version  string
	url      string
}

type LatestEntry struct {
	HistoryURL string
	Title      string
	KB         string
	URL        string
}

func (f *Fetcher) FetchLatestForBuild(ctx context.Context, build int) (*LatestEntry, error) {
	return f.FetchLatestForVersion(ctx, "", build)
}

func (f *Fetcher) FetchLatestForVersion(ctx context.Context, version string, build int) (*LatestEntry, error) {
	historyURL := historyURLByVersion(version, build)
	if historyURL == "" {
		return nil, fmt.Errorf("unsupported windows version %q build %d", version, build)
	}
	return f.FetchLatestFromHistory(ctx, historyURL)
}

func historyURLByVersion(version string, build int) string {
	version = strings.TrimSpace(strings.ToUpper(version))

	if version != "" {
		for _, page := range historyPages {
			if page.version == version {
				return page.url
			}
		}
	}

	for _, page := range historyPages {
		if build >= page.minBuild {
			return page.url
		}
	}

	return ""
}

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
