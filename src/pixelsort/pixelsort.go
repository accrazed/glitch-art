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
	// create slice channel
	sliceChan := getSlices(ps.image)

	// break slices into sortable chunks
	chunkChan := getChunks(sliceChan, ps.ThresholdFunc)

	// sort chunk channel
	sortChunks(chunkChan, ps.SorterFunc, ps.invert)

	return ps.image
}

func getSlices(image *image.RGBA64) (slices chan []uint8) {
	min, max := image.Bounds().Min, image.Bounds().Max

	slices = make(chan []uint8)
	go func(slices chan []uint8) {
		defer close(slices)
		for y := min.Y; y < max.Y; y++ {
			slices <- image.Pix[image.PixOffset(0, y):image.PixOffset(0, y+1)]
		}
	}(slices)

	return slices
}

func getChunks(slices <-chan []uint8, thresholdFunc ThresholdFunc) (chunks chan []uint8) {
	chunks = make(chan []uint8)

	go func(chunks chan<- []uint8, slices <-chan []uint8) {
		defer close(chunks)

		var wg sync.WaitGroup
		for slice := range slices {
			wg.Add(1)
			go func(chunks chan<- []uint8, slice []uint8) {
				defer wg.Done()

				start := 0
				for end := 0; end < len(slice); end += 8 {
					pixel := slice[end : end+8 : end+8]
					if !thresholdFunc(pixel) {
						continue
					}

					if end-start == 0 {
						continue
					}
					if end-start == 8 {
						start = end
						continue
					}

					chunks <- slice[start:end]
					start = end
				}
				if len(slice)-start > 8 {
					chunks <- slice[start:]
				}
			}(chunks, slice)
		}
		wg.Wait()
	}(chunks, slices)

	return chunks
}

func sortChunks(chunks <-chan []uint8, sorterFunc SorterFunc, invert bool) {
	var wg sync.WaitGroup

	for chunk := range chunks {
		wg.Add(1)

		go func(chunk []uint8) {
			defer wg.Done()

			cs := &chunkSorter{
				pixels:   chunk,
				sortFunc: sorterFunc,
				invert:   invert,
			}
			sort.Sort(cs)
		}(chunk)
	}
	wg.Wait()
}

type chunkSorter struct {
	pixels   []uint8
	sortFunc SorterFunc
	invert   bool
}

func (p *chunkSorter) Len() int {
	return len(p.pixels) / 8
}

func (p *chunkSorter) Swap(i, j int) {
	tmp := make([]uint8, 8)

	iPix, jPix := p.pixels[i*8:i*8+8], p.pixels[j*8:j*8+8]
	copy(tmp, iPix)
	copy(iPix, jPix)
	copy(jPix, tmp)
}

func (p *chunkSorter) Less(i, j int) bool {
	iPix, jPix := p.pixels[i*8:i*8+8], p.pixels[j*8:j*8+8]

	return p.sortFunc(iPix, jPix) != p.invert
}

var _ sort.Interface = (*chunkSorter)(nil)
