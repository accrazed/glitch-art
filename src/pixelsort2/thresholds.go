package pixelsort

type ThresholdFunc func(pixel []uint8) bool

const ThresholdScale = 100

var ThresholdTypes = []string{"LoThresholdColorMean", "HiThresholdColorMean", "OutThresholdColorMean", "InThresholdColorMean"}

// Only colors below the threshold will sort
func (ps *PixelSort) LoThresholdColorMean(pixel []uint8) bool {
	r := uint32(pixel[0])<<8 | uint32(pixel[1])
	g := uint32(pixel[2])<<8 | uint32(pixel[3])
	b := uint32(pixel[4])<<8 | uint32(pixel[5])

	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v > ps.threshold
}

// Only colors above the threshold will sort
func (ps *PixelSort) HiThresholdColorMean(pixel []uint8) bool {
	r := uint32(pixel[0])<<8 | uint32(pixel[1])
	g := uint32(pixel[2])<<8 | uint32(pixel[3])
	b := uint32(pixel[4])<<8 | uint32(pixel[5])

	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v < ps.threshold
}

// Only colors inside the range will sort
func (ps *PixelSort) OutThresholdColorMean(pixel []uint8) bool {
	r := uint32(pixel[0])<<8 | uint32(pixel[1])
	g := uint32(pixel[2])<<8 | uint32(pixel[3])
	b := uint32(pixel[4])<<8 | uint32(pixel[5])

	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v < ps.threshold || v > ThresholdScale-ps.threshold
}

// Only colors outside the range will sort
func (ps *PixelSort) InThresholdColorMean(pixel []uint8) bool {
	r := uint32(pixel[0])<<8 | uint32(pixel[1])
	g := uint32(pixel[2])<<8 | uint32(pixel[3])
	b := uint32(pixel[4])<<8 | uint32(pixel[5])

	v := int(float64((r+g+b)/3) / 0xFFFF * ThresholdScale)

	return v > ps.threshold && v < ThresholdScale-ps.threshold
}
