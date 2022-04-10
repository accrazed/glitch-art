package pixelsort

import (
	"image"
	"os"
)

type NewOpt func(*PixelSort) *PixelSort

type SortDir int64

const (
	Horizontal SortDir = iota
	Vertical
)

type PixelSort struct {
	seed int64

	image   *image.Image
	format  string
	sortDir SortDir
}

func Must(ps *PixelSort, err error) *PixelSort {
	if err != nil {
		panic(err)
	}
	return ps
}

func New(path string, opts ...NewOpt) (*PixelSort, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	image, format, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	ps := &PixelSort{
		image:  &image,
		format: format,
	}

	for _, opt := range opts {
		ps = opt(ps)
	}

	return ps, nil
}

func WithSortDir(dir SortDir) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.sortDir = dir
		return ps
	}
}

func WithSeed(seed int64) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.seed = seed
		return ps
	}
}

func (ps *PixelSort) Sort() *image.Image {
	return ps.image
}
