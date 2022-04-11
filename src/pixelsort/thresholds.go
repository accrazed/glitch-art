package pixelsort

type ThresholdFunc func(r, g, b, a uint32) bool

const ThresholdScale = 256

func (ps *PixelSort) ThresholdColorMean(r, g, b, a uint32) bool {
	return int(float64((r+g+b)/3)/0xFFFF*ThresholdScale) < ps.threshold
}
