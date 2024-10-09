package prompt

import (
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/eunomie/docker-runx/internal/pizza"
	"github.com/eunomie/docker-runx/runkit"
)

func SelectAction(actions []runkit.Action) string {
	var (
		action string
		form   = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select the action to run").
					Options(pizza.Map[runkit.Action, huh.Option[string]](actions, func(action runkit.Action) huh.Option[string] {
						return huh.NewOption(action.ID+envStr(action.Env), action.ID)
					})...).
					Value(&action),
			),
		)
		err = form.Run()
	)

	if err != nil {
		return ""
	}

	return action
}

func envStr(env []string) string {
	if len(env) == 0 {
		return ""
	}
	return " (required env: " + strings.Join(env, ", ") + ")"
}
