package pixelsort

import (
	"image"
	"math/rand"
	"sort"
	"sync"

	"github.com/accrazed/glitch-art/src/lib"
)

type PixelSort struct {
	image      *image.RGBA64
	direction  lib.Direction
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

	ps := &PixelSort{
		image:     image,
		threshold: -1,
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

func (ps *PixelSort) Sort() *image.RGBA64 {
	min, max := ps.image.Bounds().Min, ps.image.Bounds().Max

	// create slice channel
	slices := make(chan []uint8)
	go func(slices chan []uint8) {
		defer close(slices)
		for y := min.Y; y < max.Y; y++ {
			slices <- ps.image.Pix[ps.image.PixOffset(0, y):ps.image.PixOffset(0, y+1)]
		}
	}(slices)

	// break slices into sortable chunks
	chunks := make(chan [][]uint8)
	go func(chunks chan<- [][]uint8, slices <-chan []uint8) {
		defer close(chunks)

		var wg sync.WaitGroup
		for slice := range slices {
			wg.Add(1)
			go func(chunks chan<- [][]uint8, slice []uint8) {
				defer wg.Done()

				chunk := [][]uint8{}
				for pix := 0; pix < len(slice); pix += 8 {
					pixel := slice[pix : pix+8 : pix+8]
					if !ps.ThresholdFunc(pixel) {
						chunk = append(chunk, pixel)
						continue
					}

					if len(chunk) == 0 {
						continue
					}
					if len(chunk) == 1 {
						chunk = (chunk)[1:]
						continue
					}

					chunks <- chunk
					chunk = [][]uint8{}
				}
				if len(chunk) > 1 {
					chunks <- chunk
				}
			}(chunks, slice)
		}
		wg.Wait()
	}(chunks, slices)

	// sort chunks
	var wg sync.WaitGroup
	for chunk := range chunks {
		wg.Add(1)
		go func(chunk [][]uint8) {
			defer wg.Done()

			refs := make([][]uint8, len(chunk))
			copy(refs, chunk)
			sort.Slice(chunk, func(i, j int) bool {
				return ps.SorterFunc(chunk[i], chunk[j]) != ps.invert
			})

			for i, pixel := range chunk {
				// pixel = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
				refs[i] = pixel
			}
		}(chunk)
	}
	wg.Wait()

	return ps.image
}

type pixSorter struct {
	src      []uint8
	pixels   [][]uint8
	sortFunc SorterFunc
}

func (p *pixSorter) Len() int {
	return len(p.pixels)
}

func (p *pixSorter) Swap(i, j int) {
	p.pixels[i], p.pixels[j] = p.pixels[j], p.pixels[i]

	tmp := make([]uint8, 8)
	copy(tmp, p.pixels[i])
	copy(p.pixels[i], p.pixels[j])
	copy(p.pixels[j], tmp)
}

func (p *pixSorter) Less(i, j int) bool {
	return p.sortFunc(p.pixels[i], p.pixels[j])
}

var _ sort.Interface = (*pixSorter)(nil)
