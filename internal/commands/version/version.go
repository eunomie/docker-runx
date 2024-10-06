package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/constants"
)

func NewCmd(_ command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show Docker RunX version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(constants.Short())
		},
	}

	return cmd
}
