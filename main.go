package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sunshineplan/imgconv"
	"golang.org/x/term"
)

const CHARACTERS = " .:coPO?@■"
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
			buf.WriteByte(CHARACTERS[idx])
			buf.WriteByte(CHARACTERS[idx])
		}

		buf.WriteByte('\n')
	}

	return buf.String(), nil
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
	square := flag.Bool("square", false, "Use square aspect ratio")

	flag.Parse()

	if *square {
		*width = *height
	}

	ext := strings.Split(*input, ".")[1]

	switch ext {
	case "jpg", "jpeg", "png":
		buf, err := imageToAscii(*input, *width, *height)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf)
	default:
		fmt.Printf("[-] Invalid format: .%s\n", ext)
	}
}
