package decorate

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/runkit"
)

var (
	configFile string
	config     []byte
	readmeFile string
	readme     []byte
	tag        string
	err        error
)

func NewCmd(_ command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decorate OPTIONS IMAGE",
		Short: "Decorate an image by attaching a runx manifest",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if tag == "" {
				return errors.New("missing required flag: --tag")
			}

			if configFile == "" && readmeFile == "" {
				return errors.New("provide at least one of the following flags: --with-config, --with-readme")
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

			fmt.Println("decorating...")
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Start()

			err = runkit.Decorate(cmd.Context(), image, tag, config, readme)
			s.Stop()
			if err != nil {
				return err
			}
			fmt.Printf("successfully pushed image %s decorated with runx configuration %s to %s\n", image, configFile, tag)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&configFile, "with-config", "runx.yaml", "Path to the runx manifest file")
	f.StringVar(&readmeFile, "with-readme", "README.md", "Path to the README file")
	f.StringVarP(&tag, "tag", "t", "", "Tag to push the decorated image to")

	return cmd
}
