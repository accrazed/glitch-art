package smear

import (
	"image"
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

type Smearer struct {
	img *image.RGBA64
	r   *rand.Rand

	strength uint

	smearPos int
	smearLen int
}

func New(path string, opts ...SmearOpt) (*Smearer, error) {
	img, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}

	s := &Smearer{
		img: img,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func (s *Smearer) Smear() image.Image {
	if s.smearPos == 0 {
		s.smearPos = s.r.Intn(s.img.Bounds().Dy())
	}
	if s.smearLen == 0 {
		s.smearLen = s.r.Intn(s.img.Bounds().Dy() / int(s.strength))
	}

	for i := 0; i < s.smearLen; i++ {
		slStart := s.img.PixOffset(0, s.smearPos)
		slEnd := s.img.PixOffset(0, s.smearPos+1)
		s.img.Pix = append(
			s.img.Pix[:slEnd], append(s.img.Pix[slStart:slEnd], s.img.Pix[slEnd:]...)...)
	}

	s.img.Rect = image.Rect(
		s.img.Rect.Min.X, s.img.Rect.Min.Y,
		s.img.Rect.Max.X, s.img.Rect.Max.Y+s.smearLen)

	return s.img
}
