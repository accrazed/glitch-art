package pixelsort

import (
	"image"
	"image/color"
	"os"
	"sort"
)

type NewOpt func(*PixelSort) *PixelSort

type SortDir int64

const (
	Horizontal SortDir = iota
	Vertical
)

type PixelSort struct {
	seed            int64
	image           *image.RGBA64
	format          string
	sortDir         SortDir
	threshold       int
	breaksThreshold ThresholdFunc
	pixelSorterFunc SortFunc
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

	roImage, format, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	image := copyImage(roImage)

	ps := &PixelSort{
		image:           image,
		format:          format,
		threshold:       -1,
		pixelSorterFunc: MeanComp,
	}
	ps.breaksThreshold = ps.ThresholdColorMean

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

func (ps *PixelSort) Sort() image.Image {
	// TODO: swap based on ps.sortDir
	min, max := ps.image.Bounds().Min, ps.image.Bounds().Max

	ch := make(chan bool)
	for row := min.X; row < max.X; row++ {
		go func(row int) {
			for col := min.Y; col < max.Y; col++ {
				chunk := ps.getChunk(row, col)
				sort.Slice(chunk, func(i, j int) bool {
					return ps.pixelSorterFunc(chunk[i], chunk[j])
				})
				for i, c := range chunk {
					ps.image.Set(row, col+i, c)
				}
				col += len(chunk)
			}
			ch <- true
		}(row)
	}

	for i := min.X; i < max.X; i++ {
		<-ch
	}

	return ps.image
}

// getChunkLength returns a chunk of pixels in the range from (row,col) according to ps.compFunc
func (ps *PixelSort) getChunk(row, col int) []color.Color {
	var cur int
	res := make([]color.Color, 0)

	for cur = col; cur < ps.image.Bounds().Max.Y; cur++ {
		r, g, b, a := ps.image.At(row, cur).RGBA()
		if ps.breaksThreshold(r, g, b, a) {
			break
		}
		res = append(res, ps.image.At(row, cur))
	}

	return res
}

func copyImage(img image.Image) *image.RGBA64 {
	res := image.NewRGBA64(img.Bounds())

	min, max := img.Bounds().Min, img.Bounds().Max
	for row := min.X; row < max.X; row++ {
		for col := min.Y; col < max.Y; col++ {
			res.Set(row, col, img.At(row, col))
		}
	}

	return res
}
