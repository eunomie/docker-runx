package decorate

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/runkit"
)

var (
	noConfigFile bool
	configFile   string
	config       []byte
	noReadmeFile bool
	readmeFile   string
	readme       []byte
	tag          string
	err          error
)

func NewCmd(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decorate OPTIONS IMAGE",
		Short: "Decorate an image by attaching a runx manifest",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if tag == "" {
				return errors.New("missing required flag: --tag")
			}

			if noConfigFile {
				configFile = ""
			}
			if noReadmeFile {
				readmeFile = ""
			}

			if configFile == "" && readmeFile == "" {
				return errors.New("you need to provide a config file or a readme file, or both")
			}
			if configFile != "" {
				config, err = os.ReadFile(configFile)
				if err != nil {
					return err
				}
			}

			if readmeFile != "" {
				readme, err = os.ReadFile(readmeFile)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			image := args[0]

			err = spinner.New().
				Type(spinner.Globe).
				Title(" Decorating and pushing...").
				Action(func() {
					err = runkit.Decorate(cmd.Context(), image, tag, config, readme)
					if err != nil {
						_, _ = fmt.Fprintln(dockerCli.Err(), err)
						os.Exit(1)
					}
				}).
				Run()
			if err != nil {
				return err
			}
			fmt.Printf("successfully pushed image %s decorated with runx configuration %s to %s\n", image, configFile, tag)
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&noConfigFile, "no-config", false, "Do not attach a runx configuration to the image")
	f.StringVar(&configFile, "with-config", "runx.yaml", "Path to the runx configuration file")
	f.BoolVar(&noReadmeFile, "no-readme", false, "Do not attach a README to the image")
	f.StringVar(&readmeFile, "with-readme", "README.md", "Path to the README file")
	f.StringVarP(&tag, "tag", "t", "", "Tag to push the decorated image to")

	return cmd
}
