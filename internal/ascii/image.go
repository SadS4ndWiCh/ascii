package ascii

import (
	"image"
	"os"

	"github.com/sunshineplan/imgconv"
)

func openImage(src string) (image.Image, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func FromImage(src string, w int, h int) (*ASCII, error) {
	img, err := openImage(src)
	if err != nil {
		return nil, err
	}

	img = imgconv.Resize(img, &imgconv.ResizeOption{Width: w, Height: h})

	return New(img), nil
}
