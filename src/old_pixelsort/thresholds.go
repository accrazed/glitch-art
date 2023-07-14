package pixelsort

import "image/color"

type ThresholdFunc func(color color.Color) bool

const ThresholdScale = 100

var ThresholdTypes = []string{"LoThresholdColorMean", "HiThresholdColorMean", "OutThresholdColorMean", "InThresholdColorMean"}

// Only colors below the threshold will sort
func (ps *PixelSort) LoThresholdColorMean(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v > ps.threshold
}

// Only colors above the threshold will sort
func (ps *PixelSort) HiThresholdColorMean(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v < ps.threshold
}

// Only colors inside the range will sort
func (ps *PixelSort) OutThresholdColorMean(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v < ps.threshold || v > ThresholdScale-ps.threshold
}

// Only colors outside the range will sort
func (ps *PixelSort) InThresholdColorMean(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v > ps.threshold && v < ThresholdScale-ps.threshold
}
