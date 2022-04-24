package pixelsort

import (
	"image"
	"image/color"
	"sort"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*PixelSort) *PixelSort

type PixelSort struct {
	image *image.RGBA64

	seed          int64
	direction     lib.Direction
	invert        bool
	threshold     int
	chunkLimit    int
	ThresholdFunc ThresholdFunc
	SorterFunc    SorterFunc
}

func Must(ps *PixelSort, err error) *PixelSort {
	if err != nil {
		panic(err)
	}
	return ps
}

func New(path string, opts ...NewOpt) (*PixelSort, error) {
	image, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}

	ps := &PixelSort{
		image:     image,
		threshold: -1,
	}
	ps.SorterFunc = ps.MeanComp
	ps.ThresholdFunc = ps.OutThresholdColorMean

	for _, opt := range opts {
		ps = opt(ps)
	}

	if ps.chunkLimit == 0 {
		if ps.image.Rect.Dx() > ps.image.Rect.Dy() {
			ps.chunkLimit = ps.image.Rect.Dx()
		}
		ps.chunkLimit = ps.image.Rect.Dy()
	}

	if ps.threshold == -1 {
		ps.threshold = int((ps.seed % ThresholdScale))
	}

	return ps, nil
}

func (ps *PixelSort) Sort() image.Image {
	min, max := ps.image.Bounds().Min, ps.image.Bounds().Max
	pMin, pMax, sMin, sMax := min.X, max.X, min.Y, max.Y
	if ps.direction == lib.Horizontal {
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
					return ps.SorterFunc(chunk[i], chunk[j]) != ps.invert
				})

				// Save data
				for i, c := range chunk {
					sl, p := slice, pos+i
					if ps.direction == lib.Horizontal {
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

	for c, lim := pos, 0; c < sMax && lim < ps.chunkLimit; c, lim = c+1, lim+1 {
		sl := slice
		cur := c
		if ps.direction == lib.Horizontal {
			sl, cur = cur, sl
		}

		if ps.ThresholdFunc(ps.image.At(sl, cur)) {
			break
		}
		res = append(res, ps.image.At(sl, cur))
	}

	return res
}
