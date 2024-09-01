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

	"github.com/SadS4ndWiCh/ascii/internal/bytes"
)

type playArgs struct {
	input string
}

type PlayCommand struct {
	fs   *flag.FlagSet
	args playArgs
}

func NewPlayCommand() *PlayCommand {
	cmd := &PlayCommand{
		fs: flag.NewFlagSet("play", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.args.input, "i", "./input.mp4", "The input file")

	return cmd
}

func (cmd *PlayCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *PlayCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *PlayCommand) Run() error {
	file, err := os.Open(cmd.args.input)
	if err != nil {
		return err
	}
	defer file.Close()

	asciiReader := bytes.NewReader(file)

	signature, err := asciiReader.ReadBytes(5)
	if err != nil {
		return err
	}

	if string(signature) != "ascii" {
		return errors.New("invalid ascii file")
	}

	width, _ := asciiReader.ReadInt16()
	height, _ := asciiReader.ReadInt16()
	asciiReader.ReadInt8() // ignore aspect

	length := width*height + height
	buf := make([]byte, length)

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
