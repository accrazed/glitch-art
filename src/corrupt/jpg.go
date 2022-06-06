package corrupt

import (
	"math/rand"
)

const (
	FF               = 255
	DA               = 218
	JPEGMarkerOffset = 2
)

type JPEGBytes []byte

type Corrupter interface {
	Corrupt(jpeg JPEGBytes) (JPEGBytes, error)
}
type JPEGCorrupt struct {
	// settings
}

func (c *JPEGCorrupt) Corrupt(jpeg JPEGBytes) (JPEGBytes, error) {
	// find data start
	var start int
	for i := 0; i < len(jpeg); i++ {
		if jpeg[i] == FF && jpeg[i+1] == DA {
			StartOfScanLen := int(jpeg[i+2])<<8 + int(jpeg[i+3])
			start = i + StartOfScanLen
			break
		}
	}

	header := jpeg[:start]
	data := jpeg[start : len(jpeg)-JPEGMarkerOffset]
	eof := jpeg[len(jpeg)-JPEGMarkerOffset:]

	// corrupt
	corrupted := make([]byte, 0)
	for _, b := range data {
		//? REMOVING THIS IS WHY CORRUPTING JPGS SOMETIMES CREATES IMAGES THAT END ABRUBTLY?
		bump := byte(0)
		if rand.Intn(100000) < 1 {
			bump = byte(rand.Intn((1 << 32)))
		}
		corrupted = append(corrupted, b+bump)
	}

	res := append(append(header, corrupted...), eof...)
	return res, nil
}

var _ Corrupter = (*JPEGCorrupt)(nil)
