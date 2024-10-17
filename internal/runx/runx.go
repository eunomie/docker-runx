package runx

import (
	"context"
	"errors"

	"github.com/charmbracelet/huh/spinner"

	"github.com/eunomie/docker-runx/internal/tui"
	"github.com/eunomie/docker-runx/runkit"
)

func SelectAction(action, src, defaultAction string) string {
	if action != "" {
		return action
	}

	if conf, ok := runkit.GetLocalConfig().Image(src); ok && conf.Default != "" {
		return conf.Default
	}

	return defaultAction
}

func Get(ctx context.Context, fd uintptr, localCache runkit.Cache, src string) (*runkit.RunKit, error) {
	var (
		err error
		rk  *runkit.RunKit
	)

	if tui.IsATTY(fd) {
		var getErr error
		err = spinner.New().
			Type(spinner.Globe).
			Title(" Fetching runx details...").
			Action(func() {
				rk, getErr = runkit.Get(ctx, localCache, src)
			}).Run()
		err = errors.Join(err, getErr)
	} else {
		rk, err = runkit.Get(ctx, localCache, src)
	}

	return rk, err
}
