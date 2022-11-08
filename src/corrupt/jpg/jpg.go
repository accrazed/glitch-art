//! I don't know how to write this in a good way ;-;
package jpg

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/accrazed/glitch-art/src/corrupt"
)

const (
	FF               = 255
	DA               = 218
	JPEGMarkerOffset = 2
)

type JPEGCorrupt struct {
	r *rand.Rand
	// Recommended 1e4. The lower the strength, the more likely the a byte is to corrupt
	strength int

	header, data, eof []byte
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

	var start int
	for i := 0; i < len(bb); i++ {
		if bb[i] == FF && bb[i+1] == DA {
			SOSLen := int(bb[i+2])<<8 + int(bb[i+3])
			start = i + SOSLen
			break
		}
	}

	c := &JPEGCorrupt{
		strength: -1,
		header:   bb[:start],
		data:     bb[start : len(bb)-JPEGMarkerOffset],
		eof:      bb[len(bb)-JPEGMarkerOffset:],
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.r == nil {
		c.r = rand.New(rand.NewSource(time.Now().Unix()))
	}

	return c, nil
}

func (c *JPEGCorrupt) Build() []byte {
	return append(c.header, append(c.data, c.eof...)...)
}

func (c *JPEGCorrupt) Corrupt() *JPEGCorrupt {
	// corrupt
	c.data = corrupt.New(c.data).
		SetRand(c.r).
		SetStrength(c.strength).
		Delete().Replace().Defect().
		Data()

	return c
}
