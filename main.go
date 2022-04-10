package main

import (
	"image/png"
	"os"
	"time"

	ps "github.com/accrazed/glitch-art/src/pixelsort"
)

func main() {
	pixSort := ps.Must(
		ps.New(
			"snaaaake.png",
			ps.WithSortDir(ps.Horizontal),
			ps.WithSeed(time.Now().Unix()),
		),
	)

	img := pixSort.Sort()

	f, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, *img)
	if err != nil {
		panic(err)
	}
}
