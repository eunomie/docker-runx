package cache

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/constants"
)

func NewCmd(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: fmt.Sprintf("Manage %s cache and temporary files", constants.FullProductName),
	}
	cmd.AddCommand(
		dfNewCmd(dockerCli),
		pruneNewCmd(dockerCli),
	)

	return cmd
}
