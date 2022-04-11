package pixelsort

import "image/color"

func SaturationComp(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()
	_ = (aR + aG + aB) / 3
	_ = (bR + bG + bB) / 3
	return aR < bR
}
