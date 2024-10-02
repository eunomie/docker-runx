package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	clidocstool "github.com/docker/cli-docs-tool"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/commands/root"
)

const defaultSourcePath = "docs/reference/"

type options struct {
	source  string
	formats []string
}

func gen(opts *options) error {
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return err
	}
	cmd := &cobra.Command{
		Use:               "docker [OPTIONS] COMMAND [ARG...]",
		Short:             "The base command for the Docker CLI.",
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(root.NewCmd(dockerCLI, true))

	c, err := clidocstool.New(clidocstool.Options{
		Root:      cmd,
		SourceDir: opts.source,
		Plugin:    true,
	})
	if err != nil {
		return err
	}

	for _, format := range opts.formats {
		switch format {
		case "md":
			if err = c.GenMarkdownTree(cmd); err != nil {
				return err
			}
		case "yaml":
			if err = c.GenYamlTree(cmd); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown format %q", format)
		}
	}

	return nil
}

func run() error {
	opts := &options{}
	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	flags.StringVar(&opts.source, "source", defaultSourcePath, "Docs source folder")
	flags.StringSliceVar(&opts.formats, "formats", []string{}, "Format (md, yaml)")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	if len(opts.formats) == 0 {
		return fmt.Errorf("docs format required")
	}
	return gen(opts)
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
