package jpg

import "math/rand"

type NewOpt func(*JPEGCorrupt)

func WithSeed(seed int64) NewOpt {
	return func(jc *JPEGCorrupt) {
		jc.r = rand.New(rand.NewSource(seed))
	}
}

func WithCorruptStrength(corruptStrength int) NewOpt {
	return func(jc *JPEGCorrupt) {
		jc.strength = corruptStrength
	}
}
