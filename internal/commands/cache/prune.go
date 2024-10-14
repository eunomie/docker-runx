package cache

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/runkit"
	"github.com/spf13/cobra"
)

var (
	force bool
)

func pruneNewCmd(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove all cache entries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cache := runkit.NewLocalCache(dockerCli)

			if !force {
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Are you sure you want to remove all cache entries?").
							Value(&force))).Run()
				if err != nil {
					return err
				}
			}

			if !force {
				_, _ = fmt.Fprintln(dockerCli.Out(), "Cancelled")
				return nil
			}

			err := cache.Erase()
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(dockerCli.Out(), "Cached data deleted")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Do not prompt for confirmation")

	return cmd
}
