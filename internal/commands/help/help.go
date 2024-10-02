package help

import (
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
)

const (
	commandName = "help"
)

func NewCmd(_ command.Cli, rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: "Display information about the available commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.Help()
		},
	}

	return cmd
}
