package channelshift

import (
	"image"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*ChannelShift) *ChannelShift

type ChannelShift struct {
	format string
	tran   Translate
	image  *image.RGBA64
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
		tran: Translate{
			r: image.Point{100, 0},
			g: image.Point{30, 0},
			b: image.Point{500, 0},
		},
	}
	for _, opt := range opts {
		cs = opt(cs)
	}

	return cs, nil
}

func RedShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.tran.r.X = x
		cs.tran.r.Y = y
		return cs
	}
}

func GreenShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.tran.g.X = x
		cs.tran.g.Y = y
		return cs
	}
}

func BlueShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.tran.b.X = x
		cs.tran.b.Y = y
		return cs
	}
}

func AlphaShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.tran.a.X = x
		cs.tran.a.Y = y
		return cs
	}
}

func (cs *ChannelShift) Shift() image.Image {
	width, height := cs.image.Rect.Dx(), cs.image.Rect.Dy()

	img := lib.CopyImage(cs.image)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			old := RGBA64toPix(x, y, cs.image.Stride)

			newRed := RGBA64toPix((x+cs.tran.r.X)%width, (y+cs.tran.r.Y)%height, cs.image.Stride)
			newGreen := RGBA64toPix((x+cs.tran.g.X)%width, (y+cs.tran.g.Y)%height, cs.image.Stride)
			newBlue := RGBA64toPix((x+cs.tran.b.X)%width, (y+cs.tran.b.Y)%height, cs.image.Stride)
			newAlpha := RGBA64toPix((x+cs.tran.a.X)%width, (y+cs.tran.a.Y)%height, cs.image.Stride)

			// Red
			img.Pix[old+0] = cs.image.Pix[newRed+0]
			img.Pix[old+1] = cs.image.Pix[newRed+1]
			// Green
			img.Pix[old+2] = cs.image.Pix[newGreen+2]
			img.Pix[old+3] = cs.image.Pix[newGreen+3]
			// Blue
			img.Pix[old+4] = cs.image.Pix[newBlue+4]
			img.Pix[old+5] = cs.image.Pix[newBlue+5]
			// Alpha
			img.Pix[old+6] = cs.image.Pix[newAlpha+6]
			img.Pix[old+7] = cs.image.Pix[newAlpha+7]
		}
	}

	return img
}

func RGBA64toPix(x, y, stride int) int {
	return y*stride + x*8
}
