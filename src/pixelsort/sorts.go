package pixelsort

import "image/color"

type SorterFunc func(a, b color.Color) bool

func MeanComp(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()
	aVal := (aR + aG + aB) / 3
	bVal := (bR + bG + bB) / 3
	return aVal < bVal
}

func RedComp(a, b color.Color) bool {
	aR, _, _, _ := a.RGBA()
	bR, _, _, _ := b.RGBA()
	return aR < bR
}

func GreenComp(a, b color.Color) bool {
	_, aG, _, _ := a.RGBA()
	_, bR, _, _ := b.RGBA()
	return aG < bR
}

func BlueComp(a, b color.Color) bool {
	_, _, aB, _ := a.RGBA()
	_, _, bB, _ := b.RGBA()
	return aB < bB
}
