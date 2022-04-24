package pixelsort

import (
	"fmt"
	"reflect"

	"github.com/accrazed/glitch-art/src/lib"
)

func WithDirection(dir lib.Direction) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.direction = dir
		return ps
	}
}

func WithSeed(seed int64) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.seed = seed
		return ps
	}
}

func WithThreshold(threshold int) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.threshold = threshold
		return ps
	}
}

func WithChunkLimit(lim int) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.chunkLimit = lim
		return ps
	}
}

func WithSortFuncString(sortFunc string) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		vPS := reflect.ValueOf(ps)

		vMethod := vPS.MethodByName(sortFunc)
		zero := reflect.Value{}
		if vMethod == zero {
			panic(fmt.Sprintf("sort func %s not found", sortFunc))
		}
		vPS.Elem().FieldByName("SorterFunc").Set(vMethod)

		return ps
	}
}

func WithThresholdFuncString(thresholdFunc string) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		vPS := reflect.ValueOf(ps)

		vMethod := vPS.MethodByName(thresholdFunc)
		zero := reflect.Value{}
		if vMethod == zero {
			panic(fmt.Sprintf("threshold func %s not found", thresholdFunc))
		}
		vPS.Elem().FieldByName("ThresholdFunc").Set(vMethod)

		return ps
	}
}

func WithInvert(invert bool) NewOpt {
	return func(ps *PixelSort) *PixelSort {
		ps.invert = invert
		return ps
	}
}
