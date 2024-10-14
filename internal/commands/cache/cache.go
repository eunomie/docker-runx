package cache

import (
	"fmt"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/constants"
	"github.com/spf13/cobra"
)

func NewCmd(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: fmt.Sprintf("Manage %s cache and temporary files", constants.FullProductName),
	}
	cmd.AddCommand(
		dfNewCmd(dockerCli),
		//pruneNewCmd(dockerCli),
	)

	return cmd
}
