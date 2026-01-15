package updates

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Fetcher struct {
	Client *http.Client
	UA     string
}

func NewFetcher() *Fetcher {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   8 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   8 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       60 * time.Second,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   10,
	}

	return &Fetcher{
		Client: &http.Client{
			Timeout:   20 * time.Second,
			Transport: transport,
		},
		UA: "WinUpdBrief/0.1 (+https://example.invalid) Go-http-client",
	}
}

func (f *Fetcher) GetHTML(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	// 简单指数退避
	backoff := []time.Duration{0, 700 * time.Millisecond, 1500 * time.Millisecond}

	for attempt := 0; attempt < len(backoff); attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff[attempt]):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		// 固定 en-us
		req.Header.Set("User-Agent", f.UA)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip")

		resp, err := f.Client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		// 429/5xx 重试
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("http %d", resp.StatusCode)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("http %d", resp.StatusCode)
		}

		var body io.Reader = resp.Body
		if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return nil, err
			}
			defer gr.Close()
			body = gr
		}

		// 4MB 上限
		b, err := io.ReadAll(io.LimitReader(body, 4<<20))
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	if lastErr == nil {
		lastErr = errors.New("request failed")
	}
	return nil, lastErr
}
