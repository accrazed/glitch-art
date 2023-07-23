package smear

import (
	"image"
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

type Smearer struct {
	img *image.RGBA64
	r   *rand.Rand
}

func New(path string, seed int64) (*Smearer, error) {
	img, err := lib.NewImage(path)
	if err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(seed))

	s := &Smearer{
		img: img,
		r:   r,
	}

	return s, nil
}

func (s *Smearer) Smear() image.Image {
	smearPos := s.r.Intn(s.img.Bounds().Dy())
	smearLen := s.r.Intn(s.img.Bounds().Dy()/10) + (s.img.Bounds().Dy() / 10)

	for i := 0; i < smearLen; i++ {
		slStart := s.img.PixOffset(0, smearPos)
		slEnd := s.img.PixOffset(0, smearPos+1)
		s.img.Pix = append(
			s.img.Pix[:slEnd], append(s.img.Pix[slStart:slEnd], s.img.Pix[slEnd:]...)...)
	}

	s.img.Rect = image.Rect(
		s.img.Rect.Min.X, s.img.Rect.Min.Y,
		s.img.Rect.Max.X, s.img.Rect.Max.Y+smearLen)

	return s.img
}
