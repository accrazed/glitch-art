package channelshift

import (
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*ChannelShift)

func WithChunks(dist int) NewOpt {
	return func(cs *ChannelShift) {
		cs.chunk = dist
	}
}

func WithSeed(seed int64) NewOpt {
	return func(cs *ChannelShift) {
		cs.rand = rand.New(rand.NewSource(seed))
	}
}

func WithDirection(direction lib.Direction) NewOpt {
	return func(cs *ChannelShift) {
		cs.direction = direction
	}
}

func WithOffsetVolatility(offsetVol int) NewOpt {
	return func(cs *ChannelShift) {
		cs.offsetVol = offsetVol
	}
}

func WithChunkVolatility(chunkVol int) NewOpt {
	return func(cs *ChannelShift) {
		cs.chunkVol = chunkVol
	}
}

func WithAnimate(animate int) NewOpt {
	return func(cs *ChannelShift) {
		cs.animate = animate
	}
}

func WithRedShift(x, y int) NewOpt {
	return func(cs *ChannelShift) {
		cs.translate.r.X = x
		cs.translate.r.Y = y
	}
}

func WithGreenShift(x, y int) NewOpt {
	return func(cs *ChannelShift) {
		cs.translate.g.X = x
		cs.translate.g.Y = y
	}
}

func WithBlueShift(x, y int) NewOpt {
	return func(cs *ChannelShift) {
		cs.translate.b.X = x
		cs.translate.b.Y = y
	}
}

func WithAlphaShift(x, y int) NewOpt {
	return func(cs *ChannelShift) {
		cs.translate.a.X = x
		cs.translate.a.Y = y
	}
}
