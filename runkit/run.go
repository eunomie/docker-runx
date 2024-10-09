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
		command string
		args    string
		data    TemplateData
		Action  Action
	}

	TemplateData struct {
		Ref  string
		Env  map[string]string
		Opts map[string]string
	}
)

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

	runnable := Runnable{
		Action:  action,
		command: "docker run",
		data: TemplateData{
			Ref: ref,
			Env: map[string]string{},
		},
	}

	for _, env := range action.Env {
		if v, ok := os.LookupEnv(env); !ok {
			return nil, fmt.Errorf("environment variable %q is required", env)
		} else {
			runnable.data.Env[env] = v
		}
	}

	return &runnable, nil
}

func (r *Runnable) compute() error {
	tmpl, err := template.New("runx").Funcs(template.FuncMap{
		"env": func(envName string) string {
			return r.data.Env[envName]
		},
		"opt": func(optName string) string {
			return r.data.Opts[optName]
		},
	}).Parse(r.Action.Command)
	if err != nil {
		return err
	}

	out := strings.Builder{}
	err = tmpl.Execute(&out, r.data)
	if err != nil {
		return err
	}
	r.args = out.String()

	return nil
}

func (r *Runnable) SetOptionValues(opts map[string]string) error {
	r.data.Opts = opts

	if err := r.compute(); err != nil {
		return err
	}

	r.Command = fmt.Sprintf("%s %s", r.command, r.args)
	return nil
}

func (r *Runnable) Run(ctx context.Context) error {
	if r.Command == "" {
		return fmt.Errorf("command not set")
	}

	parsedCmd, err := syntax.NewParser().Parse(strings.NewReader(r.Command), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(interp.Env(expand.ListEnviron(os.Environ()...)), interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}
	return runner.Run(ctx, parsedCmd)
}
