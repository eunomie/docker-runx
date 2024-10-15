package prompt

import (
	"cmp"
	"errors"
	"strconv"
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
		err       error
		form      *huh.Form
		fields    []huh.Field
		asked     []string
		boolAsked []string
	)

	for _, opt := range action.Options {
		if _, ok := opts[opt.Name]; ok {
			continue
		}
		if opt.NoPrompt {
			continue
		}
		opt := opt

		var (
			title       = cmp.Or(opt.Prompt, cmp.Or(opt.Description, opt.Name))
			description = sugar.If(title != opt.Description, opt.Description, "")
		)
		switch opt.Type {
		case runkit.OptTypeInput:
			fields = append(fields,
				huh.NewInput().
					Title(title).
					Key(opt.Name).
					Description(description).
					Placeholder(opt.Default).
					Suggestions(sugar.If(opt.Default != "", []string{opt.Default}, nil)).
					Validate(checkRequired(opt.Required)))
			asked = append(asked, opt.Name)
		case runkit.OptTypeSelect:
			fields = append(fields,
				huh.NewSelect[string]().
					Title(title).
					Key(opt.Name).
					Description(description).
					Validate(checkRequired(opt.Required)).
					Options(pizza.Map(opt.Values, func(str string) huh.Option[string] {
						return huh.NewOption(str, str).Selected(str == opt.Default)
					})...))
			asked = append(asked, opt.Name)
		case runkit.OptTypeConfirm:
			fields = append(fields,
				huh.NewConfirm().
					Title(title).
					Key(opt.Name).
					Description(description))
			boolAsked = append(boolAsked, opt.Name)
		}
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
	for _, optName := range boolAsked {
		opts[optName] = strconv.FormatBool(form.GetBool(optName))
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
