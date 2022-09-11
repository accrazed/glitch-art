package pixelsort

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sort"
	"sync"

	"github.com/accrazed/glitch-art/src/lib"
)

type PixelSort struct {
	image      *image.RGBA64
	direction  lib.Direction
	mask       [][]bool
	invert     bool
	chunkLimit int
	r          *rand.Rand

	ThresholdFunc ThresholdFunc
	threshold     int

	SorterFunc SorterFunc
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

	mask := make([][]bool, image.Rect.Dy())
	for i := range mask {
		mask[i] = make([]bool, image.Rect.Dx())
	}
	ps := &PixelSort{
		image:     image,
		threshold: -1,
		mask:      mask,
	}
	ps.ThresholdFunc = ps.OutThresholdColorMean
	ps.SorterFunc = ps.MeanComp

	for _, opt := range opts {
		opt(ps)
	}

	if ps.chunkLimit == 0 {
		if ps.image.Rect.Dx() > ps.image.Rect.Dy() {
			ps.chunkLimit = ps.image.Rect.Dx()
		}
		ps.chunkLimit = ps.image.Rect.Dy()
	}

	if ps.threshold == -1 {
		ps.threshold = ps.r.Intn(ThresholdScale)
	}

	return ps, nil
}

func (ps *PixelSort) Image() image.Image {
	return ps.image
}

func (ps *PixelSort) Sort() *PixelSort {
	min, max := ps.image.Bounds().Min, ps.image.Bounds().Max
	pMin, pMax, sMin, sMax := min.X, max.X, min.Y, max.Y
	if ps.direction == lib.Horizontal {
		pMin, pMax, sMin, sMax = min.Y, max.Y, min.X, max.X
	}

	ps.processThresholdMask()

	// Iterate through each slice of pixels
	wg := sync.WaitGroup{}
	for slice := pMin; slice < pMax; slice++ {
		wg.Add(1)
		go func(slice int) {
			defer wg.Done()
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
		}(slice)
	}
	wg.Wait()

	return ps
}

// getChunkLength returns a chunk of pixels in the range from (slice,pos) according to the threshold mask
func (ps *PixelSort) getChunk(slice, pos, slMax int) []color.Color {
	res := make([]color.Color, 0)

	for c, lim := pos, 0; c < slMax && lim < ps.chunkLimit; c, lim = c+1, lim+1 {
		sl := slice
		cur := c
		if ps.direction == lib.Horizontal {
			sl, cur = cur, sl
		}

		if ps.checkPixel(sl, cur) {
			break
		}
		res = append(res, ps.image.At(sl, cur))
	}

	return res
}

// checkPixel refers to the threshMask to see if the current pixel passed a threshold and should be considered a "break" for the pixel sort
func (ps *PixelSort) checkPixel(x, y int) bool {
	return ps.mask[y][x]
}

// processThresholdMask runs ps.ThresholdFunc on every pixel in an image, updating the threshMask as it processes
func (ps *PixelSort) processThresholdMask() error {
	if ps.ThresholdFunc == nil {
		return fmt.Errorf("processThresholdMask: ps.ThresholdFunc is nil")
	}

	for x := 0; x < ps.image.Rect.Dx(); x++ {
		for y := 0; y < ps.image.Rect.Dy(); y++ {
			passes := ps.ThresholdFunc(ps.image.At(x, y))
			if passes {
				ps.mask[y][x] = true
			}
		}
	}

	return nil
}
