package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SadS4ndWiCh/ascii/internal/ascii"
	"golang.org/x/term"
)

type imageArgs struct {
	input  string
	width  int
	height int
	aspect string
}

type ImageCommand struct {
	fs   *flag.FlagSet
	args imageArgs
}

func NewImageCommand() *ImageCommand {
	cmd := &ImageCommand{
		fs: flag.NewFlagSet("image", flag.ContinueOnError),
	}

	defaultWidth, defaultHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		defaultWidth = 166
		defaultHeight = 37
	}

	cmd.fs.StringVar(&cmd.args.input, "i", "./input.jpg", "The input file")
	cmd.fs.IntVar(&cmd.args.width, "w", defaultWidth, "Set a custom width")
	cmd.fs.IntVar(&cmd.args.height, "h", defaultHeight, "Set a custom height")
	cmd.fs.StringVar(&cmd.args.aspect, "a", "default", "The aspect ratio (s/square, p/portrait)")

	return cmd
}

func (cmd *ImageCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *ImageCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *ImageCommand) Run() error {
	switch cmd.args.aspect {
	case "s", "square":
		cmd.args.width = cmd.args.height * 3
	case "p", "portrait":
		cmd.args.width = int(float64(cmd.args.height) * 1.5)
	}

	ext := filepath.Ext(cmd.args.input)
	if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
		return fmt.Errorf("file format '%s' not supported", ext)
	}

	image, err := ascii.FromImage(cmd.args.input, cmd.args.width, cmd.args.height)
	if err != nil {
		return err
	}

	buf, err := image.ToASCII()
	if err != nil {
		return err
	}

	fmt.Println(buf)

	return nil
}
