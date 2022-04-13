package channelshift

import (
	"image"
	"math"
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*ChannelShift) *ChannelShift

type ChannelShift struct {
	format     string
	translate  Translate
	image      *image.RGBA64
	seed       int64
	rand       *rand.Rand
	volatility int
	chunk      int
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
	img, format, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}

	cs := &ChannelShift{
		image:  img,
		format: format,
	}
	for _, opt := range opts {
		cs = opt(cs)
	}

	if cs.rand == nil {
		cs.rand = rand.New(rand.NewSource(0))
	}

	return cs, nil
}

func WithChunks(dist int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.chunk = dist
		return cs
	}
}

func WithSeed(seed int64) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.seed = seed
		cs.rand = rand.New(rand.NewSource(seed))
		return cs
	}
}

func WithVolatility(volatility int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.volatility = volatility
		return cs
	}
}
func RedShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.r.X = x
		cs.translate.r.Y = y
		return cs
	}
}

func GreenShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.g.X = x
		cs.translate.g.Y = y
		return cs
	}
}

func BlueShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.b.X = x
		cs.translate.b.Y = y
		return cs
	}
}

func AlphaShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.a.X = x
		cs.translate.a.Y = y
		return cs
	}
}

func (cs *ChannelShift) Shift() image.Image {
	w, h := cs.image.Rect.Dx(), cs.image.Rect.Dy()

	outImg := lib.CopyImage(cs.image)
	offsetIndex := -1
	offset := 0

	ch := make(chan bool)
	for x := 0; x < w; x++ {
		if cs.chunk != 0 && x/cs.chunk != offsetIndex {
			offsetIndex = x / cs.chunk
			offset = (cs.rand.Int() % (cs.volatility * 2)) - cs.volatility
		}

		go func(x, offset int) {
			for y := 0; y < h; y++ {
				old := lib.RGBA64toPix(x, y, cs.image.Stride)

				newR := lib.RGBA64toPix(
					(x+cs.translate.r.X+offset)%w,
					(y+cs.translate.r.Y+offset)%h,
					cs.image.Stride)
				newG := lib.RGBA64toPix(
					(x+cs.translate.g.X+offset)%w,
					(y+cs.translate.g.Y+offset)%h,
					cs.image.Stride)
				newB := lib.RGBA64toPix(
					(x+cs.translate.b.X+offset)%w,
					(y+cs.translate.b.Y+offset)%h,
					cs.image.Stride)
				newA := lib.RGBA64toPix(
					(x+cs.translate.a.X+offset)%w,
					(y+cs.translate.a.Y+offset)%h,
					cs.image.Stride)

				newR = int(math.Abs(float64(newR)))
				newG = int(math.Abs(float64(newG)))
				newB = int(math.Abs(float64(newB)))
				newA = int(math.Abs(float64(newA)))

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

			}
			ch <- true
		}(x, offset)
	}
	for i := 0; i < w; i++ {
		<-ch
	}

	return outImg
}
