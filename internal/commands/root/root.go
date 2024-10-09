package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/commands/decorate"
	"github.com/eunomie/docker-runx/internal/commands/help"
	"github.com/eunomie/docker-runx/internal/commands/version"
	"github.com/eunomie/docker-runx/internal/constants"
	"github.com/eunomie/docker-runx/internal/tui"
	"github.com/eunomie/docker-runx/runkit"
)

var docs bool

func NewCmd(dockerCli command.Cli, isPlugin bool) *cobra.Command {
	var (
		name = commandName(isPlugin)
		cmd  = &cobra.Command{
			Use:   fmt.Sprintf("%s [IMAGE] [ACTION]", name),
			Short: "Docker Run, better",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 || len(args) > 2 {
					return cmd.Help()
				}

				src := args[0]
				rk, err := runkit.Get(cmd.Context(), src)
				if err != nil {
					return err
				}

				if docs {
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(rk.Readme))
					return nil
				}

				if len(args) == 2 {
					action := args[1]
					runnable, err := rk.GetRunnable(action)
					if err != nil {
						return err
					}

					_, _ = fmt.Fprintf(dockerCli.Out(), "\nRunning the following command:\n$ %s\n\n", runnable)

					return runnable.Run(cmd.Context())
				}

				return cmd.Help()
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
		decorate.NewCmd(dockerCli),
	)

	f := cmd.Flags()
	f.BoolVarP(&docs, "docs", "d", false, "Print the documentation of the image")

	return cmd
}

func commandName(isPlugin bool) string {
	name := constants.SubCommandName
	if !isPlugin {
		name = constants.BinaryName
	}
	return name
}
