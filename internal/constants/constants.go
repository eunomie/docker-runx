package constants

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	ProductName    = "RunX"
	SubCommandName = "runx"

	FullProductName = "Docker " + ProductName

	BinaryName    = "docker-" + SubCommandName
	PluginCommand = "docker " + SubCommandName
)

var (
	Version    string = "(devel)"
	Revision   string
	LastCommit time.Time
	DirtyBuild bool
	UserAgent  string
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
			case "vcs.time":
				LastCommit, _ = time.Parse(time.RFC3339, setting.Value)
			case "vcs.modified":
				DirtyBuild = setting.Value == "true"
			}
		}
	}

	UserAgent = fmt.Sprintf("%s/%s go/%s git-commit/%s", BinaryName, Version, runtime.Version(), Revision)
}

func Short() string {
	parts := make([]string, 0, 3)
	parts = append(parts, Version)
	if Revision != "unknown" && Revision != "" {
		commit := Revision
		if len(commit) > 7 {
			commit = commit[:7]
		}
		parts = append(parts, commit)
		if DirtyBuild {
			parts = append(parts, "dirty")
		}
	}
	return strings.Join(parts, "-")
}
