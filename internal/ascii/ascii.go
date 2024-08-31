package ascii

import (
	"image"
	"image/color"
	"strings"
)

const DEFAULT_CHARACTERS = " .:-=+*oO#%8@■"

// Luminance method
func luminance(r, g, b uint8) float64 {
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}

type ASCIIABLE interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

type ASCII struct {
	r          ASCIIABLE
	characters []rune
}

func New(r ASCIIABLE) *ASCII {
	return &ASCII{r, []rune(DEFAULT_CHARACTERS)}
}

func NewWithCharacters(r ASCIIABLE, characters []rune) *ASCII {
	return &ASCII{r, characters}
}

func (a *ASCII) ToASCII() (string, error) {
	bounds := a.r.Bounds()
	buf := strings.Builder{}
	charactersLen := len(a.characters)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := a.r.At(x, y).RGBA()

			intensity := luminance(uint8(r/256), uint8(g/256), uint8(b/256))
			idx := int(intensity / 255 * float64(charactersLen-1))

			buf.WriteRune(a.characters[idx])
		}

		buf.WriteByte('\n')
	}

	return buf.String(), nil
}
