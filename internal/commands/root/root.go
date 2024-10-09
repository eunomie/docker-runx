package root

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/commands/decorate"
	"github.com/eunomie/docker-runx/internal/commands/help"
	"github.com/eunomie/docker-runx/internal/commands/version"
	"github.com/eunomie/docker-runx/internal/constants"
	"github.com/eunomie/docker-runx/internal/prompt"
	"github.com/eunomie/docker-runx/internal/registry"
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
				var (
					src    string
					action string
					lc     = runkit.GetLocalConfig()
				)

				switch len(args) {
				case 0:
					src = lc.Ref
					if src == "" {
						return cmd.Help()
					}
				case 1:
					if lc.Ref == "" {
						src = args[0]
					} else {
						// here we need to know if the argument is an image or an action
						// there's no easy way, so what we'll do is to check if the argument is a reachable image
						if registry.ImageExist(cmd.Context(), args[0]) {
							// the image exist, let's say we override the default reference
							src = args[0]
						} else {
							// we can't access the image, let's say it's an action
							src = lc.Ref
							action = args[0]
						}
					}
				case 2:
					src = args[0]
					action = args[1]
				default:
					return cmd.Help()
				}

				var (
					err error
					rk  *runkit.RunKit
				)

				err = spinner.New().
					Type(spinner.Globe).
					Title(" Fetching runx details...").
					Action(func() {
						rk, err = runkit.Get(cmd.Context(), src)
						if err != nil {
							_, _ = fmt.Fprintln(dockerCli.Err(), err)
							os.Exit(1)
						}
					}).Run()
				if err != nil {
					return err
				}

				if docs {
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(rk.Readme+"\n---\n"+mdActions(rk)))
					return nil
				}

				if list || action == "" {
					if tui.IsATTY(dockerCli.In().FD()) {
						action := prompt.SelectAction(rk.Config.Actions)
						if action != "" {
							return run(cmd.Context(), dockerCli.Err(), rk, action)
						}
					} else {
						_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(mdActions(rk)))
					}
					return nil
				}

				if action != "" {
					return run(cmd.Context(), dockerCli.Err(), rk, action)
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

func run(ctx context.Context, out io.Writer, rk *runkit.RunKit, action string) error {
	runnable, err := rk.GetRunnable(action)
	if err != nil {
		return err
	}

	opts, err := prompt.Ask(runnable.Action)
	if err != nil {
		return err
	}

	if err = runnable.SetOptionValues(opts); err != nil {
		return err
	}

	_, _ = fmt.Fprintln(out, tui.Markdown(fmt.Sprintf(`
> **Running the following command:**
>
>     %s
`, runnable.Command)))

	return runnable.Run(ctx)
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
		if action.Desc != "" {
			s.WriteString(fmt.Sprintf("  - `%s`: %s\n", action.ID, action.Desc))
		} else {
			s.WriteString(fmt.Sprintf("  - `%s`\n", action.ID))
		}
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
