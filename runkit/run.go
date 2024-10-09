package runkit

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type (
	Runnable struct {
		Command string
		Args    string
	}

	TemplateData struct {
		Ref string
	}
)

func GetRunnable(actions *RunKit, action string) (*Runnable, error) {
	return actions.GetRunnable(action)
}

func (actions *RunKit) GetRunnable(action string) (*Runnable, error) {
	for _, a := range actions.Actions {
		if a.ID == action {
			return a.GetRunnable(actions.src)
		}
	}

	return nil, fmt.Errorf("action %s not found", action)
}

func (action Action) GetRunnable(ref string) (*Runnable, error) {
	if action.Type != ActionTypeRun {
		return nil, fmt.Errorf("unsupported action type %s", action.Type)
	}

	runnable := Runnable{
		Command: "docker run",
	}

	tmpl, err := template.New(action.ID).Parse(action.Command)
	if err != nil {
		return nil, err
	}

	out := strings.Builder{}
	err = tmpl.Execute(&out, TemplateData{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}
	runnable.Args = out.String()

	return &runnable, nil
}

func (r Runnable) String() string {
	return fmt.Sprintf("%s %s", r.Command, r.Args)
}

func (r Runnable) Run(ctx context.Context) error {
	parsedCmd, err := syntax.NewParser().Parse(strings.NewReader(r.String()), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(interp.Env(expand.ListEnviron(os.Environ()...)), interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}
	return runner.Run(ctx, parsedCmd)
}
