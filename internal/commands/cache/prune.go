package cache

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/sugar"
	"github.com/eunomie/docker-runx/runkit"
)

var (
	force bool
	all   bool
)

func pruneNewCmd(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove cache entries not accessed recently",
		Long:  "By default remove cache entries not accessed in the last 30 days. Use --all/-a to remove all cache entries.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			cache := runkit.NewLocalCache(dockerCli)

			if !force {
				err = huh.NewConfirm().
					Title(sugar.If(all, "Are you sure you want to remove all cache entries?", "Are you sure you want to remove cache entries not accessed in the last 30 days?")).
					Value(&force).Run()
				if err != nil {
					return err
				}
			}

			if !force {
				_, _ = fmt.Fprintln(dockerCli.Out(), "Cancelled")
				return nil
			}

			if !all {
				err = cache.EraseNotAccessedInLast30Days()
			} else {
				err = cache.EraseAll()
			}
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(dockerCli.Out(), "Cached data deleted")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Do not prompt for confirmation")
	flags.BoolVarP(&all, "all", "a", false, "Remove all cache entries")

	return cmd
}
