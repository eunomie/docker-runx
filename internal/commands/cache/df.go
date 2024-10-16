package cache

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/formatter/tabwriter"
	"github.com/eunomie/docker-runx/runkit"
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

			w := tabwriter.NewWriter(&str, 0, 0, 1, ' ', 0)
			_, _ = fmt.Fprintln(w, "Digest\tSize\tLast Access")
			for _, e := range entries {
				t := "--"
				if e.LastAccess != nil {
					t = e.LastAccess.Format("2006-01-02 15:04:05")
				}
				_, _ = fmt.Fprintf(w,
					"%s\t%s\t%s\n",
					e.Digest,
					humanize.Bytes(uint64(e.Size)),
					t)
			}
			_ = w.Flush()
			str.WriteString(fmt.Sprintf("Total: %s\n", humanize.Bytes(uint64(totalSize))))

			_, _ = fmt.Fprintln(dockerCli.Out(), str.String())

			return nil
		},
	}
}
