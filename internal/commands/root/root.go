package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/commands/help"
	"github.com/eunomie/docker-runx/internal/commands/version"
	"github.com/eunomie/docker-runx/internal/constants"
)

func NewCmd(dockerCli command.Cli, isPlugin bool) *cobra.Command {
	var (
		name = commandName(isPlugin)
		cmd  = &cobra.Command{
			Use:   fmt.Sprintf("%s [command]", name),
			Short: "Docker Run, better",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return cmd.Help()
				}
				_ = cmd.Help()
				return cli.StatusError{
					StatusCode: 1,
					Status:     fmt.Sprintf("unknown command: %q", args[0]),
				}
			},
		}
	)

	if isPlugin {
		originalPreRunE := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if originalPreRunE != nil {
				if err := originalPreRunE(cmd, args); err != nil {
					return err
				}
			}
			return nil
		}
	} else {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		cmd.TraverseChildren = true
		cmd.DisableFlagsInUseLine = true
		cli.DisableFlagsInUseLine(cmd)
	}

	cmd.AddCommand(
		help.NewCmd(dockerCli, cmd),
		version.NewCmd(dockerCli),
	)

	return cmd
}

func commandName(isPlugin bool) string {
	name := constants.SubCommandName
	if !isPlugin {
		name = constants.BinaryName
	}
	return name
}
