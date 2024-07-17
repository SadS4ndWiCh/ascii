package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/imgconv"
	"golang.org/x/term"
)

// const CHARACTERS = " .:coPO?@■"

const CHARACTERS = " .:-=+*oO#%8@■"
const CHARACTERS_LENGTH = len(CHARACTERS)

func rgbtogray(r, g, b float64) float64 {
	return 0.3*r + 0.59*g + 0.11*b
}

func imageToAscii(src string, w int, h int) (string, error) {
	img, err := imgconv.Open(src)
	if err != nil {
		return "", err
	}

	buf := strings.Builder{}
	img = imgconv.Resize(img, &imgconv.ResizeOption{Width: w, Height: h})
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r = r % 256
			g = g % 256
			b = b % 256

			intensity := rgbtogray(float64(r), float64(g), float64(b))
			idx := int(intensity / 255 * float64(CHARACTERS_LENGTH-1))

			buf.WriteByte(CHARACTERS[idx])
		}

		buf.WriteByte('\n')
	}

	return buf.String(), nil
}

func videoToAscii(src string, w int, h int) error {
	if err := os.Mkdir("ascii-tmp", 0764); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	splitInChunksCmd := exec.Command(
		"ffmpeg.exe",
		"-i",
		src,
		"-r",
		"30",
		"./ascii-tmp/%05d.jpg",
	)

	_, err := splitInChunksCmd.Output()
	if err != nil {
		return err
	}

	e, err := os.ReadDir("./ascii-tmp")
	if err != nil {
		return err
	}

	frames := len(e)

	name := filepath.Base(src)
	file, err := os.Create(fmt.Sprintf("%s.ascii", name))
	if err != nil {
		return err
	}
	defer file.Close()

	for frame := range frames {
		frameSrc := fmt.Sprintf("./ascii-tmp/%05d.jpg", frame)

		buf, err := imageToAscii(frameSrc, w, h)
		if err != nil {
			continue
		}

		file.Write([]byte(buf))
	}

	if err := os.RemoveAll("./ascii-tmp"); err != nil {
		return err
	}

	return nil
}

func playVideo(src string, w int, h int) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	offset := (w * h) + h
	buf := make([]byte, offset)

	for {
		n, err := file.Read(buf)
		if err != nil {
			return err
		} else if n == 0 {
			break
		}

		fmt.Printf("\033[0;0H%s", buf)
		time.Sleep(time.Second / 30)
	}

	return nil
}

func main() {
	defaultWidth, defaultHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		defaultWidth = 166
		defaultHeight = 37
	}

	input := flag.String("i", "./input.jpg", "The input file")
	width := flag.Int("w", defaultWidth, "Set a custom width")
	height := flag.Int("h", defaultHeight, "Set a custom height")
	square := flag.Bool("s", false, "Use square aspect ratio")

	flag.Parse()

	if *square {
		*width = *height * 3
	}

	ext := filepath.Ext(*input)

	switch ext {
	case ".jpg", ".jpeg", ".png":
		buf, err := imageToAscii(*input, *width, *height)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf)
	case ".mp4", ".gif":
		if err := videoToAscii(*input, *width, *height); err != nil {
			log.Fatal(err)
		}
	case ".ascii":
		if err := playVideo(*input, *width, *height); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("[-] Invalid format: %s\n", ext)
	}
}
