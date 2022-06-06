package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/accrazed/glitch-art/src/corrupt"
)

func main() {
	rand.Seed(time.Now().Unix())

	c := &corrupt.JPEGCorrupt{}
	f, err := os.Open("img/in/arch.jpg")
	if err != nil {
		panic(err)
	}
	bb, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		outbuf, err := c.Corrupt(bb)
		if err != nil {
			panic(err)
		}
		outf, err := os.Create("out.jpg")
		if err != nil {
			panic(err)
		}
		_, err = outf.Write(outbuf)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	// cli.RunCLI()
}
