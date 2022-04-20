package lib

import (
	"image"

	"github.com/sunshineplan/imgconv"
)

func CopyImage(img image.Image) *image.RGBA64 {
	res := image.NewRGBA64(img.Bounds())

	min, max := img.Bounds().Min, img.Bounds().Max
	for row := min.X; row < max.X; row++ {
		for col := min.Y; col < max.Y; col++ {
			res.Set(row, col, img.At(row, col))
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

func RGBA64toPix(x, y, stride int) int {
	return y*stride + x*8
}
