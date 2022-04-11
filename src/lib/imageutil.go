package lib

import "image"

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
