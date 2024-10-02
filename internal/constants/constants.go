package constants

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

const (
	ProductName    = "RunX"
	SubCommandName = "runx"

	FullProductName = "Docker " + ProductName

	BinaryName    = "docker-" + SubCommandName
	PluginCommand = "docker " + SubCommandName
)

var (
	Version   string = "devel"
	commit    string
	UserAgent string
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
		}
	}

	UserAgent = fmt.Sprintf("%s/%s go/%s git-commit/%s", BinaryName, Version, runtime.Version(), commit)
}
