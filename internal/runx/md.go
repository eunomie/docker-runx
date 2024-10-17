package runx

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"

	"github.com/eunomie/docker-runx/internal/sugar"
	"github.com/eunomie/docker-runx/runkit"
)

func FullMD(rk *runkit.RunKit) string {
	return rk.Readme + "\n---\n" + MDActions(rk)
}

func MDAction(rk *runkit.RunKit, action string) string {
	var (
		act   runkit.Action
		found bool
	)
	for _, a := range rk.Config.Actions {
		if a.ID == action {
			found = true
			act = a
			break
		}
	}
	if !found {
		return fmt.Sprintf("> action %q not found\n\n%s", action, MDActions(rk))
	}

	s := strings.Builder{}
	if act.Desc != "" {
		s.WriteString(fmt.Sprintf("`%s`%s: %s\n", act.ID, sugar.If(act.IsDefault(), " (default)", ""), act.Desc))
	} else {
		s.WriteString(fmt.Sprintf("`%s`\n", act.ID))
	}
	if len(act.Env) > 0 {
		s.WriteString("\n- Environment " + plural("variable", len(act.Env)) + ":\n")
		for _, env := range act.Env {
			s.WriteString("    - `" + env + "`\n")
		}
	}
	if len(act.Options) > 0 {
		s.WriteString("\n- " + plural("Option", len(act.Options)) + ":\n")
		for _, opt := range act.Options {
			s.WriteString("    - `" + opt.Name + "`" + sugar.If(opt.Description != "", ": "+opt.Description, "") + "\n")
		}
	}
	if len(act.Shell) > 0 {
		s.WriteString("\n- Shell " + plural("command", len(act.Shell)) + ":\n")
		for name, cmd := range act.Shell {
			s.WriteString("    - `" + name + "`: `" + cmd + "`\n")
		}
	}
	s.WriteString("\n- " + capitalizedTypes[act.Type] + " command:\n")
	s.WriteString("```\n" + act.Command + "\n```\n")

	return s.String()
}

var capitalizedTypes = map[runkit.ActionType]string{
	runkit.ActionTypeRun:   "Run",
	runkit.ActionTypeBuild: "Build",
}

func MDActions(rk *runkit.RunKit) string {
	s := strings.Builder{}
	s.WriteString("# Available actions\n\n")
	if len(rk.Config.Actions) == 0 {
		s.WriteString("> No available action\n")
	} else {
		for _, action := range rk.Config.Actions {
			if action.Desc != "" {
				s.WriteString(fmt.Sprintf("  - `%s`%s: %s\n", action.ID, sugar.If(action.IsDefault(), "(default)", ""), action.Desc))
			} else {
				s.WriteString(fmt.Sprintf("  - `%s`\n", action.ID))
			}
		}

		s.WriteString("\n> Use `docker runx IMAGE ACTION --docs` to get more details about an action\n")
	}

	return s.String()
}

func plural(str string, n int) string {
	p := pluralize.NewClient()
	if n > 1 {
		return p.Plural(str)
	}
	return str
}
