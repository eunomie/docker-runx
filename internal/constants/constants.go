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
	Version   string = "(devel)"
	Revision  string
	UserAgent string
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Value == "" {
				continue
			}
			switch setting.Key {
			case "vcs.revision":
				Revision = setting.Value
			}
		}
	}

	UserAgent = fmt.Sprintf("%s/%s go/%s git-commit/%s", BinaryName, Version, runtime.Version(), Revision)
}

func Runtime() string {
	return fmt.Sprintf("%s - %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
