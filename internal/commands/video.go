package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/SadS4ndWiCh/ascii/internal/ascii"
	"github.com/SadS4ndWiCh/ascii/internal/bytes"
	"golang.org/x/term"
)

type videoArgs struct {
	input  string
	width  int
	height int
	aspect string
}

type VideoCommand struct {
	fs   *flag.FlagSet
	args videoArgs
}

func NewVideoCommand() *VideoCommand {
	cmd := &VideoCommand{
		fs: flag.NewFlagSet("video", flag.ContinueOnError),
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

func (cmd *VideoCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *VideoCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *VideoCommand) Run() error {
	switch cmd.args.aspect {
	case "s", "square":
		cmd.args.width = cmd.args.height * 3
	case "p", "portrait":
		cmd.args.width = int(float64(cmd.args.height) * 1.5)
	}

	tmpPath, err := os.MkdirTemp("", "ascii-temp")
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	defer os.RemoveAll(tmpPath)

	tmpFramesPath := filepath.Join(tmpPath, "%05d.jpg")
	splitInChunksCmd := exec.Command(
		"ffmpeg",
		"-i",
		cmd.args.input,
		"-r",
		"30",
		tmpFramesPath,
	)

	if _, err := splitInChunksCmd.Output(); err != nil {
		return err
	}

	files, err := os.ReadDir(tmpPath)
	if err != nil {
		return err
	}

	frames := []string{}
	for frame := range len(files) {
		frameFile := fmt.Sprintf("%05d.jpg", frame)
		framePath := filepath.Join(tmpPath, frameFile)

		frames = append(frames, framePath)
	}

	asciiFileName := filepath.Base(cmd.args.input)
	asciiFile, err := os.Create(fmt.Sprintf("%s.ascii", asciiFileName))
	if err != nil {
		return err
	}
	defer asciiFile.Close()

	/*
		`.ascii` file structure

		 9     5 4     3 2      1 0      0
		+-------+-------+--------+--------+
		| ascii | width | height | aspect |
		+-------+-------+--------+--------+
		|          ASCII FRAMES           |
		|    (width * height) + height    |
		+---------------------------------+
	*/
	asciiWriter := bytes.NewWriter(asciiFile)

	asciiWriter.WriteString("ascii")
	asciiWriter.WriteInt16(int16(cmd.args.width))
	asciiWriter.WriteInt16(int16(cmd.args.height))

	switch cmd.args.aspect {
	case "s", "square":
		asciiWriter.WriteInt8(1)
	case "p", "portrait":
		asciiWriter.WriteInt8(2)
	default:
		asciiWriter.WriteInt8(0)
	}

	for _, frameSrc := range frames {
		img, err := ascii.FromImage(frameSrc, cmd.args.width, cmd.args.height)
		if err != nil {
			continue
		}

		buf, err := img.ToASCII()
		if err != nil {
			continue
		}

		asciiWriter.WriteString(buf)
	}

	return nil
}
