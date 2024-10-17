package root

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/commands/cache"
	"github.com/eunomie/docker-runx/internal/commands/decorate"
	"github.com/eunomie/docker-runx/internal/commands/version"
	"github.com/eunomie/docker-runx/internal/constants"
	"github.com/eunomie/docker-runx/internal/prompt"
	"github.com/eunomie/docker-runx/internal/registry"
	"github.com/eunomie/docker-runx/internal/runx"
	"github.com/eunomie/docker-runx/internal/tui"
	"github.com/eunomie/docker-runx/runkit"
)

var (
	docs        bool
	list        bool
	ask         bool
	opts        []string
	noFlagCheck bool
	helpFlag    bool
)

func NewCmd(dockerCli command.Cli, isPlugin bool) *cobra.Command {
	var (
		name = commandName(isPlugin)
		cmd  = &cobra.Command{
			Use:   fmt.Sprintf("%s [IMAGE] [ACTION]", name),
			Short: "Docker Run, better",
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					lc                    = runkit.GetLocalConfig()
					localCache            = runkit.NewLocalCache(dockerCli)
					src, action, needHelp = parseArgs(cmd.Context(), args, lc)
					err                   error
					rk                    *runkit.RunKit
				)

				if needHelp {
					return cmd.Help()
				}

				_ = localCache.EraseNotAccessedInLast30Days()

				rk, err = runx.Get(cmd.Context(), dockerCli.In().FD(), localCache, src)
				if err != nil {
					return err
				}

				// in case the image only contains the readme, display it
				// in this case we ignore the other flags or action as there's no action
				// so we can't do anything else
				if len(rk.Config.Actions) == 0 {
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(rk.Readme))
					return nil
				}

				if docs {
					var md string
					if action != "" {
						md = runx.MDAction(rk, action)
					} else {
						md = runx.FullMD(rk)
					}
					_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(md))
					return nil
				}

				action = runx.SelectAction(action, src, rk.Config.Default)

				if list || action == "" {
					if tui.IsATTY(dockerCli.In().FD()) && len(rk.Config.Actions) > 0 {
						selectedAction := prompt.SelectAction(rk.Config.Actions)
						if selectedAction != "" {
							return runx.Run(cmd.Context(), dockerCli.Err(), rk, lc, runx.RunConfig{
								Src:       src,
								Action:    selectedAction,
								ForceAsk:  ask,
								NoConfirm: noFlagCheck,
								Opts:      opts,
							})
						}
					} else {
						_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(runx.MDActions(rk)))
					}
					return nil
				}

				if action != "" {
					return runx.Run(cmd.Context(), dockerCli.Err(), rk, lc, runx.RunConfig{
						Src:       src,
						Action:    action,
						ForceAsk:  ask,
						NoConfirm: noFlagCheck,
						Opts:      opts,
					})
				}

				return cmd.Help()
			},
		}
	)

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		var (
			lc         = runkit.GetLocalConfig()
			localCache = runkit.NewLocalCache(dockerCli)
			src        string
			action     string
			needHelp   bool
			err        error
			rk         *runkit.RunKit
			md         string
		)

		if isPlugin && len(args) > 0 && args[0] == constants.SubCommandName {
			args = args[1:]
		}

		if len(args) > 0 {
			for _, c := range cmd.Commands() {
				if strings.HasPrefix(c.Use, args[0]) {
					printHelp(c)
					return
				}
			}
		}

		if err := c.ParseFlags(args); err != nil {
			printHelp(c)
			return
		}

		args = c.Flags().Args()
		src, action, needHelp = parseArgs(c.Context(), args, lc)
		if needHelp {
			printHelp(c)
			return
		}

		_ = localCache.EraseNotAccessedInLast30Days()

		rk, err = runx.Get(c.Context(), dockerCli.In().FD(), localCache, src)
		if err != nil {
			_, _ = fmt.Fprintln(dockerCli.Err(), err)
			os.Exit(1)
		}

		if action != "" {
			md = runx.MDAction(rk, action)
		} else {
			md = runx.FullMD(rk)
		}
		_, _ = fmt.Fprintln(dockerCli.Out(), tui.Markdown(md))
	})

	cmd.PersistentFlags().BoolVarP(&helpFlag, "help", "h", false, "Print usage or runx image/action documentation")

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
		version.NewCmd(dockerCli),
		decorate.NewCmd(dockerCli),
		cache.NewCmd(dockerCli),
	)

	f := cmd.Flags()
	f.BoolVarP(&docs, "docs", "d", false, "Print the documentation of the image")
	_ = f.MarkDeprecated("docs", "use -h/--help instead")
	f.BoolVarP(&list, "list", "l", false, "List available actions")
	f.BoolVar(&ask, "ask", false, "Do not read local configuration option values and always ask them")
	f.StringArrayVar(&opts, "opt", nil, "Set an option value")
	f.BoolVarP(&noFlagCheck, "yes", "y", false, "Do not check flags before running the command")

	return cmd
}

func parseArgs(ctx context.Context, args []string, lc *runkit.LocalConfig) (src, action string, needHelp bool) {
	switch len(args) {
	case 0:
		src = lc.Ref
		if src == "" {
			needHelp = true
		}
	case 1:
		if lc.Ref == "" {
			src = args[0]
		} else {
			// here we need to know if the argument is an image or an action
			// there's no easy way, so what we'll do is to check if the argument is a reachable image
			if registry.ImageExist(ctx, args[0]) {
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
		needHelp = true
	}
	return
}

func printHelp(c *cobra.Command) {
	err := tmpl(c.OutOrStdout(), c.HelpTemplate(), c)
	if err != nil {
		c.PrintErrln(err)
	}
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	// t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func commandName(isPlugin bool) string {
	name := constants.SubCommandName
	if !isPlugin {
		name = constants.BinaryName
	}
	return name
}
