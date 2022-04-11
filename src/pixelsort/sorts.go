package pixelsort

import "image/color"

type SortFunc func(a, b color.Color) bool

func MeanComp(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()
	aVal := (aR + aG + aB) / 3
	bVal := (bR + bG + bB) / 3
	return aVal < bVal
}
