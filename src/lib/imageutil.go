package lib

import (
	"image"

	"github.com/sunshineplan/imgconv"
)

func CopyImage(img image.Image, flip ...bool) *image.RGBA64 {
	bounds := img.Bounds()
	if len(flip) > 0 && flip[0] {
		bounds = image.Rect(bounds.Min.Y, bounds.Min.X, bounds.Max.Y, bounds.Max.X)
	}

	res := image.NewRGBA64(bounds)

	min, max := res.Bounds().Min, res.Bounds().Max

	for row := min.X; row < max.X; row++ {
		for col := min.Y; col < max.Y; col++ {
			if len(flip) > 0 && flip[0] {
				res.Set(row, col, img.At(col, row))
			} else {
				res.Set(row, col, img.At(row, col))
			}
		}
	}

	return res
}

func NewImage(path string) (*image.RGBA64, error) {
	roImage, err := imgconv.Open(path)
	if err != nil {
		return nil, err
	}
	image := CopyImage(roImage)

	return image, nil
}

// Converts an X,Y pixel location to the relevant stride location (Pix format).
// See the img lib for more info.
func RGBA64toPix(x, y, stride int) int {
	return y*stride + x*8
}
