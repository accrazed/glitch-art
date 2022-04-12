package pixelsort

import (
	"image"
	"image/color"
	"sort"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*PixelSort) *PixelSort

type SortDir int

const (
	Vertical SortDir = iota
	Horizontal
)

type PixelSort struct {
	seed          int64
	image         *image.RGBA64
	format        string
	sortDir       SortDir
	invert        bool
	threshold     int
	thresholdFunc ThresholdFunc
	sorterFunc    SorterFunc
}

func Must(ps *PixelSort, err error) *PixelSort {
	if err != nil {
		panic(err)
	}
	return ps
}

func New(path string, opts ...NewOpt) (*PixelSort, error) {
	image, format, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}

	ps := &PixelSort{
		image:      image,
		format:     format,
		threshold:  -1,
		sorterFunc: MeanComp,
	}
	ps.thresholdFunc = ps.ThresholdColorMean

	for _, opt := range opts {
		ps = opt(ps)
	}

	if ps.threshold == -1 {
		ps.threshold = int((ps.seed % ThresholdScale))
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

func WithThreshold(threshold int) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.threshold = threshold
		return ps
	}
}

func WithInvert(invert bool) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.invert = invert
		return ps
	}
}

func (ps *PixelSort) Sort() image.Image {
	min, max := ps.image.Bounds().Min, ps.image.Bounds().Max
	pMin, pMax, sMin, sMax := min.X, max.X, min.Y, max.Y
	if ps.sortDir == Horizontal {
		pMin, pMax, sMin, sMax = min.Y, max.Y, min.X, max.X
	}

	// Iterate through each slice of pixels
	ch := make(chan bool)
	for slice := pMin; slice < pMax; slice++ {
		go func(slice int) {
			// Iterate through pixels
			for pos := sMin; pos < sMax; pos++ {
				// Group and sort chunk
				chunk := ps.getChunk(slice, pos, sMax)
				sort.Slice(chunk, func(i, j int) bool {
					return ps.sorterFunc(chunk[i], chunk[j]) != ps.invert
				})

				// Save data
				for i, c := range chunk {
					sl, p := slice, pos+i
					if ps.sortDir == Horizontal {
						sl, p = p, sl
					}
					ps.image.Set(sl, p, c)
				}

				pos += len(chunk)
			}
			ch <- true
		}(slice)
	}
	for i := pMin; i < pMax; i++ {
		<-ch
	}

	return ps.image
}

// getChunkLength returns a chunk of pixels in the range from (slice,pos) according to ps.compFunc
func (ps *PixelSort) getChunk(slice, pos, sMax int) []color.Color {
	res := make([]color.Color, 0)

	for c := pos; c < sMax; c++ {
		sl := slice
		cur := c
		if ps.sortDir == Horizontal {
			sl, cur = cur, sl
		}

		if ps.thresholdFunc(ps.image.At(sl, cur)) {
			break
		}
		res = append(res, ps.image.At(sl, cur))
	}

	return res
}
