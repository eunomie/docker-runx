package prompt

import (
	"cmp"
	"errors"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/eunomie/docker-runx/internal/pizza"
	"github.com/eunomie/docker-runx/internal/sugar"
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
						return huh.NewOption(
							sugar.If(action.Desc == "",
								action.ID,
								action.ID+": "+action.Desc,
							)+envStr(action.Env),
							action.ID)
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

func Ask(action *runkit.Action, opts map[string]string) (map[string]string, error) {
	if len(action.Options) == 0 {
		return nil, nil
	}

	var (
		err    error
		form   *huh.Form
		fields []huh.Field
		asked  []string
	)

	for _, opt := range action.Options {
		if _, ok := opts[opt.Name]; ok {
			continue
		}
		opt := opt
		if len(opt.Values) == 0 {
			fields = append(fields,
				huh.NewInput().
					Title(cmp.Or(opt.Prompt, cmp.Or(opt.Description, opt.Name))).
					Key(opt.Name).
					Validate(checkRequired(opt.Required)))
		} else {
			fields = append(fields,
				huh.NewSelect[string]().
					Title(cmp.Or(opt.Prompt, cmp.Or(opt.Description, opt.Name))).
					Key(opt.Name).
					Validate(checkRequired(opt.Required)).
					Options(pizza.Map(opt.Values, func(str string) huh.Option[string] {
						return huh.NewOption(str, str)
					})...))
		}
		asked = append(asked, opt.Name)
	}

	if len(fields) == 0 {
		return opts, nil
	}

	form = huh.NewForm(huh.NewGroup(fields...))
	if err = form.Run(); err != nil {
		return nil, err
	}

	for _, optName := range asked {
		opts[optName] = form.GetString(optName)
	}

	return opts, nil
}

func checkRequired(isRequired bool) func(string) error {
	return func(str string) error {
		if str == "" && isRequired {
			return errors.New("required")
		}
		return nil
	}
}
