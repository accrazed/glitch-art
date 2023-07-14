package pixelsort

type SorterFunc func(a, b []uint8) bool

var SortTypes = []string{"MeanComp", "RedComp", "GreenComp", "BlueComp"}

func (ps *PixelSort) MeanComp(a, b []uint8) bool {
	aR := uint16(a[0])<<8 | uint16(a[1])
	aG := uint16(a[2])<<8 | uint16(a[3])
	aB := uint16(a[4])<<8 | uint16(a[5])

	bR := uint16(b[0])<<8 | uint16(b[1])
	bG := uint16(b[2])<<8 | uint16(b[3])
	bB := uint16(b[4])<<8 | uint16(b[5])

	aVal := float32(uint32(aR)+uint32(aG)+uint32(aB)) / 3
	bVal := float32(uint32(bR)+uint32(bG)+uint32(bB)) / 3
	return aVal < bVal
}

func (ps *PixelSort) RedComp(a, b []uint8) bool {
	aR := uint16(a[0])<<8 | uint16(a[1])
	bR := uint16(b[0])<<8 | uint16(b[1])

	return aR < bR
}
func (ps *PixelSort) GreenComp(a, b []uint8) bool {
	aG := uint16(a[2])<<8 | uint16(a[3])
	bG := uint16(b[2])<<8 | uint16(b[3])

	return aG < bG
}

func (ps *PixelSort) BlueComp(a, b []uint8) bool {
	aB := uint16(a[4])<<8 | uint16(a[5])
	bB := uint16(b[4])<<8 | uint16(b[5])

	return aB < bB
}
