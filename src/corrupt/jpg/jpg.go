//! I don't know how to write this in a good way ;-;
package jpg

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

const (
	FF               = 255
	DA               = 218
	JPEGMarkerOffset = 2
)

type JPEGBytes []byte

type JPEGCorrupt struct {
	seed int64
	rand *rand.Rand
	// Recommended 1e4. The lower the strength, the more likely the a byte is to corrupt
	corruptStrength int

	jpeg JPEGBytes
}

func Must(jpeg *JPEGCorrupt, err error) *JPEGCorrupt {
	if err != nil {
		panic(err)
	}
	return jpeg
}

func New(path string, opts ...NewOpt) (*JPEGCorrupt, error) {
	pathSp := strings.Split(path, ".")
	ext := pathSp[len(pathSp)-1]
	if ext != "jpg" && ext != "jpeg" {
		return nil, fmt.Errorf("new: file %s isn't a jpg/jpeg", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bb, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	jc := &JPEGCorrupt{
		corruptStrength: -1,
		jpeg:            bb,
	}

	for _, opt := range opts {
		opt(jc)
	}

	return jc, nil
}

func (c *JPEGCorrupt) Corrupt() ([]byte, error) {
	// find data start
	var start int
	for i := 0; i < len(c.jpeg); i++ {
		if c.jpeg[i] == FF && c.jpeg[i+1] == DA {
			SOSLen := int(c.jpeg[i+2])<<8 + int(c.jpeg[i+3])
			start = i + SOSLen
			break
		}
	}

	header := c.jpeg[:start]
	data := c.jpeg[start : len(c.jpeg)-JPEGMarkerOffset]
	eof := c.jpeg[len(c.jpeg)-JPEGMarkerOffset:]

	// corrupt
	corrupted := make([]byte, 0)
	cpy := byte(0)
	for _, b := range data {
		// do a corruption
		if c.rand.Intn(c.corruptStrength) < 1 {
			switch c.rand.Intn(4) {
			case 0: // change bit
				b = byte(c.rand.Intn((1 << 32)))
			case 1: // ignore bit
				continue
			case 2: // copy bit
				cpy = b
			case 3: // paste bit
				b = cpy
			}
		}
		corrupted = append(corrupted, b)
	}

	res := append(append(header, corrupted...), eof...)
	return res, nil
}
