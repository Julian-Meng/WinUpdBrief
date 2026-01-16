package render

import (
	"fmt"
	"strings"
)

type ViewModel struct {
	OSName         string
	DisplayVersion string
	Build          string

	KBTitle string
	KB      string
	KBURL   string
	Summary string
}

func RenderText(vm ViewModel) string {
	var b strings.Builder

	b.WriteString("Windows Update Brief\n")
	b.WriteString(strings.Repeat("=", 50) + "\n")

	// System info
	b.WriteString("System\n")
	b.WriteString(strings.Repeat("-", 50) + "\n")
	b.WriteString(fmt.Sprintf("OS: %s\n", vm.OSName))
	b.WriteString(fmt.Sprintf("Version: %s\n", vm.DisplayVersion))
	b.WriteString(fmt.Sprintf("Build: %s\n\n", vm.Build))

	// Update info
	b.WriteString("Latest Update\n")
	b.WriteString(strings.Repeat("-", 50) + "\n")
	b.WriteString(vm.KBTitle + "\n")
	b.WriteString(fmt.Sprintf("KB: %s\n", vm.KB))
	b.WriteString(fmt.Sprintf("URL: %s\n\n", vm.KBURL))

	// Summary
	if vm.Summary != "" {
		b.WriteString("Summary\n")
		b.WriteString(strings.Repeat("-", 50) + "\n")
		b.WriteString(vm.Summary)
		b.WriteString("\n")
	}

	return b.String()
}
