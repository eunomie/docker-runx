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
		Env map[string]string
	}
)

func GetRunnable(rk *RunKit, action string) (*Runnable, error) {
	if rk == nil || len(rk.Config.Actions) == 0 {
		return nil, fmt.Errorf("no available configuration is nil")
	}
	return rk.GetRunnable(action)
}

func (rk *RunKit) GetRunnable(action string) (*Runnable, error) {
	for _, a := range rk.Config.Actions {
		if a.ID == action {
			return a.GetRunnable(rk.src)
		}
	}

	return nil, fmt.Errorf("action %s not found", action)
}

func (action Action) GetRunnable(ref string) (*Runnable, error) {
	if action.Type != ActionTypeRun {
		return nil, fmt.Errorf("unsupported action type %s", action.Type)
	}

	data := TemplateData{
		Ref: ref,
		Env: map[string]string{},
	}

	for _, env := range action.Env {
		if v, ok := os.LookupEnv(env); !ok {
			return nil, fmt.Errorf("environment variable %q is required", env)
		} else {
			data.Env[env] = v
		}
	}

	runnable := Runnable{
		Command: "docker run",
	}

	tmpl, err := template.New(action.ID).Funcs(template.FuncMap{
		"env": func(envName string) string {
			return data.Env[envName]
		},
	}).Parse(action.Command)
	if err != nil {
		return nil, err
	}

	out := strings.Builder{}
	err = tmpl.Execute(&out, data)
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
