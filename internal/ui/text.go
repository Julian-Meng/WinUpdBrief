package ui

import (
	"fmt"
	"strings"
)

type View struct {
	AppName        string
	OSName         string
	DisplayVersion string
	Build          string
	UpdateTitle    string
	KB             string
	KBURL          string
	Summary        string
}

const (
	colorReset = "\x1b[0m"
	colorCyan  = "\x1b[36;1m"
	colorBlue  = "\x1b[34;1m"
	colorGreen = "\x1b[32;1m"
	colorGray  = "\x1b[90m"
	colorBold  = "\x1b[1m"
)

func RenderOverview(v View) string {
	var b strings.Builder

	title := v.AppName
	if title == "" {
		title = "WinUpdBrief"
	}

	b.WriteString(colorCyan)
	b.WriteString(title)
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(colorGray)
	b.WriteString(strings.Repeat("=", 58))
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%sSystem:%s %s\n", colorBold, colorReset, v.OSName))
	b.WriteString(fmt.Sprintf("%sVersion:%s %s\n", colorBold, colorReset, v.DisplayVersion))
	b.WriteString(fmt.Sprintf("%sBuild:%s %s\n\n", colorBold, colorReset, v.Build))

	b.WriteString(colorBlue)
	b.WriteString("Latest Update")
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(colorGray)
	b.WriteString(strings.Repeat("-", 58))
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%sTitle:%s %s\n", colorBold, colorReset, v.UpdateTitle))
	b.WriteString(fmt.Sprintf("%sKB:%s %s\n", colorBold, colorReset, v.KB))
	b.WriteString(fmt.Sprintf("%sURL:%s %s\n\n", colorBold, colorReset, v.KBURL))
	b.WriteString(colorGreen)
	b.WriteString("Preview")
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(colorGray)
	b.WriteString(strings.Repeat("-", 58))
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(previewSummary(v.Summary))
	b.WriteString("\n")
	b.WriteString(colorGray)
	b.WriteString("提示: 输入 `y` 可展开完整更新内容，或使用 `--detail` 直接显示。")
	b.WriteString(colorReset)
	b.WriteString("\n")

	return b.String()
}

func RenderDetail(v View) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(colorGreen)
	b.WriteString("Detailed Update")
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(colorGray)
	b.WriteString(strings.Repeat("-", 58))
	b.WriteString(colorReset)
	b.WriteString("\n")
	b.WriteString(strings.TrimSpace(v.Summary))
	b.WriteString("\n")

	return b.String()
}

func previewSummary(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "No summary text was found."
	}

	paragraphs := strings.Split(trimmed, "\n\n")
	if len(paragraphs) > 2 {
		paragraphs = paragraphs[:2]
	}

	preview := strings.Join(paragraphs, "\n\n")
	if preview == trimmed {
		return preview
	}

	return trimRunes(preview, 900) + "\n\n[truncated]"
}

func trimRunes(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit])
}
