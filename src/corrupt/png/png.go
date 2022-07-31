package png

import (
	"fmt"
	"math/rand"
	"time"

	png "github.com/accrazed/png/src"
)

type PNGCorrupt struct {
	r     *rand.Rand
	trans *png.Transcoder
	// Recommended 1e4. The lower the strength, the more likely the a byte is to corrupt
	strength int
}

type NewOpt func(*PNGCorrupt)

func New(path string, opts ...NewOpt) (*PNGCorrupt, error) {
	trans, err := png.NewTranscoder(path)
	if err != nil {
		return nil, fmt.Errorf("unable to create transcoder: %w", err)
	}

	c := &PNGCorrupt{
		trans:    trans,
		strength: 1e4,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.r == nil {
		c.r = rand.New(rand.NewSource(time.Now().Unix()))
	}

	return c, nil
}

func (c *PNGCorrupt) Corrupt() *PNGCorrupt {
	return c
}
