package lib

import (
	"image"
	"os"
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

func NewImage(path string) (*image.RGBA64, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	roImage, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	image := CopyImage(roImage)

	return image, format, nil
}

func RGBA64toPix(x, y, stride int) int {
	return y*stride + x*8
}
