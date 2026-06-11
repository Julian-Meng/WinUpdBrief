package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"WinUpdBrief/internal/fetcher"
	"WinUpdBrief/internal/ui"
	"WinUpdBrief/internal/windows"
)

type Options struct {
	Detail bool
}

func Run(ctx context.Context, opts Options, stdin io.Reader, stdout io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}

	info, err := windows.ReadCurrentVersion()
	if err != nil {
		return fmt.Errorf("读取 Windows 版本失败: %w", err)
	}
	if info.BuildInt <= 0 {
		return errors.New("无法识别当前 Windows Build 号")
	}

	client := fetcher.NewFetcher()
	fetchCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	entry, err := client.FetchLatestForVersion(fetchCtx, info.EffectiveDisplayVersion, info.BuildInt)
	if err != nil {
		return fmt.Errorf("获取更新信息失败: %w", err)
	}

	kb, err := client.FetchKBContent(fetchCtx, entry.URL)
	if err != nil {
		return fmt.Errorf("获取 KB 内容失败: %w", err)
	}

	view := ui.View{
		AppName:        "WinUpdBrief",
		OSName:         firstNonEmpty(info.ProductName, "Windows"),
		DisplayVersion: firstNonEmpty(info.EffectiveDisplayVersion, "Unknown"),
		Build:          firstNonEmpty(info.BuildString(), "Unknown"),
		UpdateTitle:    entry.Title,
		KB:             entry.KB,
		KBURL:          entry.URL,
		Summary:        kb.Summary,
	}

	if _, err := fmt.Fprint(stdout, ui.RenderOverview(view)); err != nil {
		return fmt.Errorf("写入输出失败: %w", err)
	}

	detail := opts.Detail
	if !detail && isInteractive(stdin) {
		detail, err = promptForDetail(stdin, stdout)
		if err != nil {
			detail = false
		}
	}

	if detail {
		if _, err := fmt.Fprint(stdout, ui.RenderDetail(view)); err != nil {
			return fmt.Errorf("写入输出失败: %w", err)
		}
	}

	return nil
}

func isInteractive(stdin io.Reader) bool {
	file, ok := stdin.(*os.File)
	if !ok {
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	return stat.Mode()&os.ModeCharDevice != 0
}

func promptForDetail(stdin io.Reader, stdout io.Writer) (bool, error) {
	if _, err := fmt.Fprint(stdout, "\n是否展开详细更新信息？[y/N]: "); err != nil {
		return false, err
	}

	reader := bufio.NewReader(stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
