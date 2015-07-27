// Package cmdtmpl provides the templated command system for mcdev
package cmdtmpl

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"text/template"
)

var ErrInvalidCommand = errors.New("invalid command")

type Command struct {
	Cmd  string
	Args []*template.Template
}

func NewCommand(args []string) (*Command, error) {
	result := new(Command)

	if len(args) == 0 {
		return nil, ErrInvalidCommand
	}

	result.Cmd = args[0]
	result.Args = make([]*template.Template, len(args)-1)

	for i, arg := range args[1:] {
		t, err := template.New("arg").Parse(arg)
		if err != nil {
			return nil, err
		}
		result.Args[i] = t
	}
	return result, nil
}

func (cmd *Command) Make(ctx interface{}) (*exec.Cmd, error) {
	args := make([]string, len(cmd.Args))

	for i, t := range cmd.Args {
		var buf bytes.Buffer
		err := t.Execute(&buf, ctx)
		if err != nil {
			return nil, err
		}
		args[i] = buf.String()
	}

	proc := exec.Command(cmd.Cmd, args...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	return proc, nil
}

func (cmd *Command) Run(ctx interface{}) error {
	proc, err := cmd.Make(ctx)
	if err != nil {
		return err
	}
	return proc.Run()
}
