package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/eunomie/docker-runx/internal/commands/root"
	"github.com/eunomie/docker-runx/internal/constants"
)

func main() {
	cmd, err := command.NewDockerCli()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if plugin.RunningStandalone() {
		err = runStandalone(cmd)
	} else {
		err = runPlugin(cmd)
	}

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runStandalone(cmd *command.DockerCli) error {
	if err := cmd.Initialize(flags.NewClientOptions()); err != nil {
		return err
	}
	rootCmd := newRootCmd(false, cmd)
	return rootCmd.Execute()
}

func runPlugin(cmd *command.DockerCli) error {
	rootCmd := newRootCmd(true, cmd)
	return plugin.RunPlugin(cmd, rootCmd, manager.Metadata{
		SchemaVersion:    "0.1.0",
		Vendor:           "Docker Inc.",
		Version:          constants.Version,
		ShortDescription: constants.FullProductName,
	})
}

func newRootCmd(isPlugin bool, dockerCli command.Cli) *cobra.Command {
	return root.NewCmd(dockerCli, isPlugin)
}
