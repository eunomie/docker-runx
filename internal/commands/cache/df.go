package cache

import (
	"fmt"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/dustin/go-humanize"
	"github.com/eunomie/docker-runx/runkit"
	"github.com/spf13/cobra"
)

func dfNewCmd(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "df",
		Short: "Show disk usage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cache := runkit.NewLocalCache(dockerCli)
			cacheDir, entries, totalSize, err := cache.ListCache()
			if err != nil {
				return err
			}

			if cacheDir == "" {
				fmt.Println("No cache directory")
				return nil
			}

			str := strings.Builder{}
			str.WriteString("Cache directory: " + cacheDir + "\n")
			str.WriteString("\n")
			for _, e := range entries {
				str.WriteString(fmt.Sprintf("%s: %s\n", e.Digest, humanize.Bytes(uint64(e.Size))))
			}
			str.WriteString(fmt.Sprintf("Total: %s\n", humanize.Bytes(uint64(totalSize))))

			_, _ = fmt.Fprintln(dockerCli.Out(), str.String())

			return nil
		},
	}
}
