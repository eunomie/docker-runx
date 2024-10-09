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
		Command    string
		command    string
		args       string
		data       TemplateData
		Action     *Action
		dockerfile string
	}

	TemplateData struct {
		Ref        string
		Env        map[string]string
		Opts       map[string]string
		Dockerfile string
	}
)

var noop = func() {}

func (rk *RunKit) GetRunnable(action string) (*Runnable, func(), error) {
	for _, a := range rk.Config.Actions {
		if a.ID == action {
			return a.GetRunnable(rk.src)
		}
	}

	return nil, noop, fmt.Errorf("action %s not found", action)
}

var rootCommands = map[ActionType]string{
	ActionTypeRun:   "docker run",
	ActionTypeBuild: "docker buildx build",
}

func (action *Action) GetRunnable(ref string) (*Runnable, func(), error) {
	rootCommand, ok := rootCommands[action.Type]
	if !ok {
		return nil, noop, fmt.Errorf("unsupported action type %s", action.Type)
	}

	runnable := Runnable{
		Action:  action,
		command: rootCommand,
		data: TemplateData{
			Ref: ref,
			Env: map[string]string{},
		},
	}

	for _, env := range action.Env {
		if v, ok := os.LookupEnv(env); !ok {
			return nil, noop, fmt.Errorf("environment variable %q is required", env)
		} else {
			runnable.data.Env[env] = v
		}
	}

	if action.DockerfileContent != "" {
		f, err := os.CreateTemp("", "runx.*.Dockerfile")
		if err != nil {
			return nil, noop, err
		}
		if _, err := f.Write([]byte(action.DockerfileContent)); err != nil {
			f.Close() //nolint:errcheck
			return nil, noop, err
		}
		runnable.dockerfile = f.Name()
		runnable.data.Dockerfile = f.Name()
		if err := f.Close(); err != nil {
			return nil, noop, err
		}
	}

	return &runnable, runnable.cleanup(), nil
}

func (r *Runnable) cleanup() func() {
	return func() {
		if r.dockerfile != "" {
			_ = os.Remove(r.dockerfile)
		}
	}
}

func (r *Runnable) compute() error {
	tmpl, err := template.New("runx").Funcs(template.FuncMap{
		"env": func(envName string) string {
			return r.data.Env[envName]
		},
		"opt": func(optName string) string {
			return r.data.Opts[optName]
		},
		"sh": func(cmdName string) (string, error) {
			cmd, ok := r.Action.Shell[cmdName]
			if !ok {
				return "", fmt.Errorf("shell command %q not found", cmdName)
			}
			return sh(context.Background(), cmd)
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

func sh(ctx context.Context, cmd string) (string, error) {
	parseCmd, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return "", err
	}

	var osOut, osErr strings.Builder
	runner, err := interp.New(interp.Env(expand.ListEnviron(os.Environ()...)), interp.StdIO(nil, &osOut, &osErr))
	if err != nil {
		return "", err
	}
	if err = runner.Run(ctx, parseCmd); err != nil {
		return "", err
	}

	return strings.TrimSpace(osOut.String()), nil
}
