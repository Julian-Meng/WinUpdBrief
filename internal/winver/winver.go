package winver

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

const currentVersionPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion`

type Info struct {
	// Raw registry values
	ProductName        string
	DisplayVersion     string
	ReleaseID          string
	CurrentBuild       string
	CurrentBuildNumber string
	UBR                uint64
	EditionID          string
	BuildLabEx         string

	// Derived / effective values
	EffectiveDisplayVersion string
	EffectiveBuild          string
	BuildInt                int
}

func ReadCurrentVersion() (Info, error) {
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		currentVersionPath,
		registry.QUERY_VALUE|registry.WOW64_64KEY,
	)
	if err != nil {
		return Info{}, fmt.Errorf("open registry key CurrentVersion: %w", err)
	}
	defer key.Close()

	info := Info{
		ProductName:        getStringValue(key, "ProductName"),
		DisplayVersion:     getStringValue(key, "DisplayVersion"),
		ReleaseID:          getStringValue(key, "ReleaseId"),
		CurrentBuild:       getStringValue(key, "CurrentBuild"),
		CurrentBuildNumber: getStringValue(key, "CurrentBuildNumber"),
		UBR:                getIntegerValue(key, "UBR"),
		EditionID:          getStringValue(key, "EditionID"),
		BuildLabEx:         getStringValue(key, "BuildLabEx"),
	}

	// ---- derived fields ----

	if info.DisplayVersion != "" {
		info.EffectiveDisplayVersion = info.DisplayVersion
	} else {
		info.EffectiveDisplayVersion = info.ReleaseID
	}

	if info.CurrentBuild != "" {
		info.EffectiveBuild = info.CurrentBuild
	} else {
		info.EffectiveBuild = info.CurrentBuildNumber
	}

	if info.EffectiveBuild != "" {
		if v, err := strconv.Atoi(info.EffectiveBuild); err == nil {
			info.BuildInt = v
		}
	}

	return info, nil
}

func getStringValue(key registry.Key, name string) string {
	value, _, err := key.GetStringValue(name)
	if err != nil {
		return ""
	}
	return value
}

func getIntegerValue(key registry.Key, name string) uint64 {
	value, _, err := key.GetIntegerValue(name)
	if err != nil {
		return 0
	}
	return value
}

func (i *Info) Format() string {
	output := "Windows CurrentVersion\n"
	output += "ProductName: " + i.ProductName + "\n"
	output += "DisplayVersion: " + i.EffectiveDisplayVersion + "\n"
	output += "CurrentBuild: " + i.EffectiveBuild + "\n"
	output += "UBR: " + fmt.Sprintf("%d", i.UBR) + "\n"
	output += "EditionID: " + i.EditionID + "\n"
	output += "BuildLabEx: " + i.BuildLabEx
	return output
}
