package commands

import (
	"errors"
	"fmt"
	"os"
)

type Command interface {
	Name() string
	Init(args []string) error
	Run() error
}

func Root(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a sub-command")
	}

	cmds := []Command{
		NewImageCommand(),
		NewVideoCommand(),
		NewPlayCommand(),
	}

	subcommand := os.Args[1]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
