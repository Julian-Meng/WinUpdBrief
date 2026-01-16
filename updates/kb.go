package updates

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type KBContent struct {
	KB      string
	Title   string
	Summary string
}

func (f *Fetcher) FetchKBContent(ctx context.Context, kbURL string) (*KBContent, error) {
	html, err := f.GetHTML(ctx, kbURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	// 标题（页面 <title> 或 h1）
	title := strings.TrimSpace(doc.Find("h1").First().Text())
	if title == "" {
		title = strings.TrimSpace(doc.Find("title").First().Text())
	}

	// Summary / Highlights
	var summaryParts []string

	doc.Find("p").Each(func(_ int, p *goquery.Selection) {
		text := strings.TrimSpace(p.Text())
		if text == "" {
			return
		}

		// 非常保守的过滤
		if len(text) < 40 {
			return
		}
		if strings.Contains(text, "Applies to") {
			return
		}

		summaryParts = append(summaryParts, text)
	})

	if len(summaryParts) == 0 {
		return nil, fmt.Errorf("no readable text found in KB page")
	}

	return &KBContent{
		Title:   title,
		Summary: strings.Join(summaryParts, "\n\n"),
	}, nil
}
