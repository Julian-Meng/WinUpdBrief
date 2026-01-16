package main

import (
	"context"
	"fmt"
	"time"

	"WinUpdBrief/render"
	"WinUpdBrief/updates"
	"WinUpdBrief/winver"
)

func main() {
	// 读取 Windows 版本信息
	info, err := winver.ReadCurrentVersion()
	if err != nil {
		fmt.Printf("Failed to read Windows version: %v\n", err)
		return
	}

	if info.BuildInt <= 0 {
		fmt.Println("Could not determine Windows build")
		return
	}

	// 抓取最新更新（带超时）
	fetcher := updates.NewFetcher()
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	entry, err := fetcher.FetchLatestForBuild(ctx, info.BuildInt)
	if err != nil {
		fmt.Printf("Failed to fetch update history: %v\n", err)
		return
	}

	// 抓取 KB 正文内容
	kb, err := fetcher.FetchKBContent(ctx, entry.URL)
	if err != nil {
		fmt.Printf("Failed to fetch KB content: %v\n", err)
		return
	}

	// 组装 ViewModel
	vm := render.ViewModel{
		OSName:         info.ProductName,
		DisplayVersion: info.EffectiveDisplayVersion,
		Build:          fmt.Sprintf("%s.%d", info.EffectiveBuild, info.UBR),

		KBTitle: entry.Title,
		KB:      entry.KB,
		KBURL:   entry.URL,
		Summary: kb.Summary,
	}

	// 5. 渲染输出
	fmt.Print(render.RenderText(vm))
}
