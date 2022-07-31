package channelshift

import (
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

type ChannelShift struct {
	translate Translate
	image     *image.RGBA64
	rand      *rand.Rand
	direction lib.Direction
	offsetVol int
	chunkVol  int
	chunk     int
	animate   int
}

type Translate struct {
	r image.Point
	g image.Point
	b image.Point
	a image.Point
}

func Must(ps *ChannelShift, err error) *ChannelShift {
	if err != nil {
		panic(err)
	}
	return ps
}

func New(path string, opts ...NewOpt) (*ChannelShift, error) {
	img, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}

	cs := &ChannelShift{image: img}
	for _, opt := range opts {
		opt(cs)
	}

	if cs.rand == nil {
		cs.rand = rand.New(rand.NewSource(0))
	}

	if cs.chunkVol > cs.chunk {
		cs.chunkVol = cs.chunk
	}

	if cs.chunk == 0 {
		if cs.image.Rect.Dx() > cs.image.Rect.Dy() {
			cs.chunk = cs.image.Rect.Dx()
		}
		cs.chunk = cs.image.Rect.Dy()
	}

	return cs, nil
}

func (cs *ChannelShift) Shift() image.Image {
	numSlices, numPos := cs.image.Rect.Dx(), cs.image.Rect.Dy()
	if cs.direction == lib.Horizontal {
		numSlices, numPos = numPos, numSlices
	}

	outImg := lib.CopyImage(cs.image)
	offsetR, offsetG, offsetB, offsetA := 0, 0, 0, 0

	ch := make(chan bool)

	// Iterate through slices perpendicular to chunking directions
	for curSlice := 0; curSlice < numSlices; {
		chunkSize := cs.chunk
		if cs.offsetVol > 0 {
			chunkSize = cs.chunk + cs.rand.Intn(cs.chunkVol*2) - cs.chunkVol
		}
		if cs.offsetVol > 0 {
			offsetR = (cs.rand.Int() % (cs.offsetVol * 2)) - cs.offsetVol
			offsetG = (cs.rand.Int() % (cs.offsetVol * 2)) - cs.offsetVol
			offsetB = (cs.rand.Int() % (cs.offsetVol * 2)) - cs.offsetVol
			offsetA = (cs.rand.Int() % (cs.offsetVol * 2)) - cs.offsetVol
		}

		// Shift each slice
		var cur int
		for cur = 0; cur < chunkSize && cur+curSlice < numSlices; cur++ {
			go func(slice, offsetR, offsetG, offsetB, offsetA int) {
				for pos := 0; pos < numPos; pos++ {
					sl, ps := numSlices, numPos
					if cs.direction == lib.Horizontal {
						sl, ps = ps, sl
						slice, pos = pos, slice
					}

					old := lib.RGBA64toPix(slice, pos, cs.image.Stride)

					rX, rY := (slice+cs.translate.r.X+offsetR)%sl,
						(pos+cs.translate.r.Y+offsetR)%ps
					gX, gY := (slice+cs.translate.g.X+offsetG)%sl,
						(pos+cs.translate.g.Y+offsetG)%ps
					bX, bY := (slice+cs.translate.b.X+offsetB)%sl,
						(pos+cs.translate.b.Y+offsetB)%ps
					aX, aY := (slice+cs.translate.a.X+offsetA)%sl,
						(pos+cs.translate.a.Y+offsetA)%ps

					newR := int(math.Abs(float64(lib.RGBA64toPix(rX, rY, cs.image.Stride))))
					newG := int(math.Abs(float64(lib.RGBA64toPix(gX, gY, cs.image.Stride))))
					newB := int(math.Abs(float64(lib.RGBA64toPix(bX, bY, cs.image.Stride))))
					newA := int(math.Abs(float64(lib.RGBA64toPix(aX, aY, cs.image.Stride))))

					// Red
					outImg.Pix[old+0] = cs.image.Pix[newR+0]
					outImg.Pix[old+1] = cs.image.Pix[newR+1]
					// Green
					outImg.Pix[old+2] = cs.image.Pix[newG+2]
					outImg.Pix[old+3] = cs.image.Pix[newG+3]
					// Blue
					outImg.Pix[old+4] = cs.image.Pix[newB+4]
					outImg.Pix[old+5] = cs.image.Pix[newB+5]
					// Alpha
					outImg.Pix[old+6] = cs.image.Pix[newA+6]
					outImg.Pix[old+7] = cs.image.Pix[newA+7]

					if cs.direction == lib.Horizontal {
						slice, pos = pos, slice
					}
				}
				ch <- true
			}(curSlice+cur, offsetR, offsetG, offsetB, offsetA)
		}
		curSlice += cur
	}
	for i := 0; i < numSlices; i++ {
		<-ch
	}

	return outImg
}

// ShiftIterate calls Shift multiple times, changing the base offset
func (cs *ChannelShift) ShiftIterate() []image.Image {
	res := make([]image.Image, 0)

	baseTr := cs.translate

	for i := 0; i < cs.animate; i++ {
		fmt.Printf("Generating channelshift frame %v...\n", i+1)
		res = append(res, cs.Shift())

		if cs.animate != 1 && cs.offsetVol != 0 {
			cs.translate.r.X = baseTr.r.X + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.r.Y = baseTr.r.Y + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.g.X = baseTr.g.X + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.g.Y = baseTr.g.Y + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.b.X = baseTr.b.X + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.b.Y = baseTr.b.Y + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.a.X = baseTr.a.X + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
			cs.translate.a.Y = baseTr.a.Y + rand.Int()%(cs.offsetVol*2) - cs.offsetVol
		}
	}

	return res
}
