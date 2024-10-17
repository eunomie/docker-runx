package runx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/eunomie/docker-runx/internal/pizza"
	"github.com/eunomie/docker-runx/internal/prompt"
	"github.com/eunomie/docker-runx/internal/tui"
	"github.com/eunomie/docker-runx/runkit"
)

type RunConfig struct {
	Src       string
	Action    string
	ForceAsk  bool
	NoConfirm bool
	Opts      []string
}

func Run(ctx context.Context, out io.Writer, rk *runkit.RunKit, lc *runkit.LocalConfig, runConfig RunConfig) error {
	runnable, cleanup, err := rk.GetRunnable(runConfig.Action)
	defer cleanup()
	if err != nil {
		return err
	}

	localOpts := map[string]string{}

	if !runConfig.ForceAsk {
		localOpts = getValuesLocal(runConfig.Src, runConfig.Action)

		for _, opt := range runConfig.Opts {
			if key, value, ok := strings.Cut(opt, "="); ok {
				localOpts[key] = value
			} else {
				return fmt.Errorf("invalid option value %s", opt)
			}
		}
	}

	options, err := prompt.Ask(runnable.Action, localOpts)
	if err != nil {
		return err
	}

	if err = runnable.SetOptionValues(options); err != nil {
		return err
	}

	mdCommand := fmt.Sprintf(`
> **Running the following command:**

    %s

---
`, runnable.Command)

	var flags []string
	if !runConfig.NoConfirm && !lc.AcceptTheRisk {
		flags, err = runnable.CheckFlags()
	}
	if err != nil {
		return err
	} else if len(flags) > 0 {
		_, _ = fmt.Fprintln(out, tui.Markdown(mdCommand+fmt.Sprintf(`
> **Some flags require your attention:**

%s
`, strings.Join(pizza.Map(flags, func(flag string) string {
			return fmt.Sprintf("- `%s`", flag)
		}), "\n"))))
		var cont bool
		err = huh.NewConfirm().Title("Continue?").Value(&cont).Run()
		if err != nil {
			return err
		}
		if !cont {
			return errors.New("aborted")
		}
	} else {
		_, _ = fmt.Fprintln(out, tui.Markdown(mdCommand))
	}

	return runnable.Run(ctx)
}

func getValuesLocal(src, action string) map[string]string {
	localOpts := make(map[string]string)

	lc := runkit.GetLocalConfig()
	img, ok := lc.Image(src)
	if !ok {
		return localOpts
	}

	if img.AllActions.Opts != nil {
		localOpts = img.AllActions.Opts
	}

	act, ok := img.Actions[action]
	if ok {
		for k, v := range act.Opts {
			localOpts[k] = v
		}
	}
	return localOpts
}
