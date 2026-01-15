package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"WinUpdBrief/internal/updates"
	"WinUpdBrief/internal/winver"
)

func main() {
	info, err := winver.ReadCurrentVersion()
	if err != nil {
		fmt.Printf("Failed to read CurrentVersion registry: %v\n", err)
		return
	}

	fmt.Println(info.Format())

	if info.BuildInt <= 0 {
		fmt.Println("Warning: could not determine Windows build, skip update fetch")
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 50))

	fetcher := updates.NewFetcher()
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	entry, err := fetcher.FetchLatestForBuild(ctx, info.BuildInt)
	if err != nil {
		fmt.Printf("Failed to fetch latest update: %v\n", err)
		return
	}

	fmt.Println("\nLatest Windows Update:")
	fmt.Printf("Title: %s\n", entry.Title)
	fmt.Printf("KB: %s\n", entry.KB)
	fmt.Printf("URL: %s\n", entry.URL)
}
