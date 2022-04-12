package pixelsort

import "image/color"

type ThresholdFunc func(color color.Color) bool

const ThresholdScale = 256

func (ps *PixelSort) ThresholdColorMean(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	return int(float64((r+g+b)/3)/0xFFFF*ThresholdScale) < ps.threshold
}
