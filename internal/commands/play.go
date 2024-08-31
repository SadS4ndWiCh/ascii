package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/term"
)

type playArgs struct {
	input  string
	width  int
	height int
	aspect string
}

type PlayCommand struct {
	fs   *flag.FlagSet
	args playArgs
}

func NewPlayCommand() *PlayCommand {
	cmd := &PlayCommand{
		fs: flag.NewFlagSet("play", flag.ContinueOnError),
	}

	defaultWidth, defaultHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		defaultWidth = 166
		defaultHeight = 37
	}

	cmd.fs.StringVar(&cmd.args.input, "i", "./input.mp4", "The input file")
	cmd.fs.IntVar(&cmd.args.width, "w", defaultWidth, "Set a custom width")
	cmd.fs.IntVar(&cmd.args.height, "h", defaultHeight, "Set a custom height")
	cmd.fs.StringVar(&cmd.args.aspect, "a", "default", "The aspect ratio (s/square, p/portrait)")

	return cmd
}

func (cmd *PlayCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *PlayCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *PlayCommand) Run() error {
	switch cmd.args.aspect {
	case "s", "square":
		cmd.args.width = cmd.args.height * 3
	case "p", "portrait":
		cmd.args.width = int(float64(cmd.args.height) * 1.5)
	}

	file, err := os.Open(cmd.args.input)
	if err != nil {
		return err
	}
	defer file.Close()

	offset := (cmd.args.width * cmd.args.height) + cmd.args.height
	buf := make([]byte, offset)

	for {
		if _, err := file.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		fmt.Printf("\033[0;0H%s", strings.TrimLeftFunc(string(buf), unicode.IsSpace))
		time.Sleep(time.Second / 30)
	}

	return nil
}
