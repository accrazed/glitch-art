package smear_test

import (
	"bufio"
	"bytes"
	"image/jpeg"
	"io"
	"os"
	"testing"

	"github.com/accrazed/glitch-art/src/smear"
	"github.com/stretchr/testify/assert"
)

func TestSmear(t *testing.T) {
	smearer, err := smear.New("fixtures/unsmeared.jpg",
		smear.WithSeed(0),
		smear.WithStrength(10))
	assert.NoError(t, err)

	img := smearer.Smear()

	fGot, err := os.Create("fixtures/smeared_got.jpg")
	assert.NoError(t, err)
	assert.NoError(t, jpeg.Encode(fGot, img, nil))

	fWant, err := os.Open("fixtures/smeared_want.jpg")
	assert.NoError(t, err)

	comp := sameFile(fGot, fWant)
	assert.True(t, comp)
}

func sameFile(gotFile, wantFile io.Reader) bool {
	gotReader := bufio.NewScanner(gotFile)
	wantReader := bufio.NewScanner(wantFile)

	for gotReader.Scan() {
		if wantReader.Scan() == false {
			return false
		}

		if !bytes.Equal(gotReader.Bytes(), wantReader.Bytes()) {
			return false
		}
	}

	return true
}
