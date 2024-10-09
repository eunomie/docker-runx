package root

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
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

var (
	docs bool
	list bool
)

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
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(rk.Readme+"\n---\n"+mdActions(rk)))
					return nil
				}

				if list || len(args) == 1 {
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(mdActions(rk)))
					return nil
				}

				if len(args) == 2 {
					action := args[1]
					runnable, err := rk.GetRunnable(action)
					if err != nil {
						return err
					}

					_, _ = fmt.Fprintln(dockerCli.Err(), tui.Markdown(fmt.Sprintf(`
> **Running the following command:**
>
>     %s
`, runnable)))

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
	f.BoolVarP(&list, "list", "l", false, "List available actions")

	return cmd
}

func commandName(isPlugin bool) string {
	name := constants.SubCommandName
	if !isPlugin {
		name = constants.BinaryName
	}
	return name
}

func mdActions(rk *runkit.RunKit) string {
	p := pluralize.NewClient()
	s := strings.Builder{}
	s.WriteString("# Available actions\n\n")
	for _, action := range rk.Config.Actions {
		s.WriteString(fmt.Sprintf("  - `%s`\n", action.ID))
		vars := "variable"
		if len(action.Env) > 1 {
			vars = p.Plural(vars)
		}
		if len(action.Env) > 0 {
			s.WriteString("    - Environment " + vars + ": " + strings.Join(tui.BackQuoteItems(action.Env), ", ") + "\n")
		}
	}

	return s.String()
}
