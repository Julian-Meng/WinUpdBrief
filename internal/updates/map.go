package updates

func HistoryURLByBuild(build int) string {
	// Win11 25H2: 26200+
	if build >= 26200 {
		return "Sorry, the Windows 11 Version 25H2 Update History page is not yet available."
	}
	// Win11 24H2: 26100+
	if build >= 26100 {
		return "Sorry, the Windows 11 Version 24H2 Update History page is not yet available."
	}
	// Win10：用统一 update history 页做兜底（也能拿到最新 Win10 KB）
	return "Sorry, the Windows 10 Update History page is not yet available."
}
