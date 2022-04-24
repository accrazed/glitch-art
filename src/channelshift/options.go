package channelshift

import (
	"math/rand"

	"github.com/accrazed/glitch-art/src/lib"
)

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

func WithDirection(direction lib.Direction) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.direction = direction
		return cs
	}
}

func WithOffsetVolatility(offsetVol int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.offsetVol = offsetVol
		return cs
	}
}

func WithChunkVolatility(chunkVol int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.chunkVol = chunkVol
		return cs
	}
}

func WithAnimate(animate int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.animate = animate
		return cs
	}
}

func WithRedShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.r.X = x
		cs.translate.r.Y = y
		return cs
	}
}

func WithGreenShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.g.X = x
		cs.translate.g.Y = y
		return cs
	}
}

func WithBlueShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.b.X = x
		cs.translate.b.Y = y
		return cs
	}
}

func WithAlphaShift(x, y int) NewOpt {
	return func(cs *ChannelShift) *ChannelShift {
		cs.translate.a.X = x
		cs.translate.a.Y = y
		return cs
	}
}
