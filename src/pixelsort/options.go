package pixelsort

import (
	"fmt"
	"reflect"

	"github.com/accrazed/glitch-art/src/lib"
)

type NewOpt func(*PixelSort)

func WithDirection(dir lib.Direction) NewOpt {
	return func(ps *PixelSort) {
		ps.direction = dir
	}
}

func WithSeed(seed int64) NewOpt {
	return func(ps *PixelSort) {
		ps.seed = seed
	}
}

func WithThreshold(threshold int) NewOpt {
	return func(ps *PixelSort) {
		ps.threshold = threshold
	}
}

func WithChunkLimit(lim int) NewOpt {
	return func(ps *PixelSort) {
		ps.chunkLimit = lim
	}
}

func WithSortFuncString(sortFunc string) NewOpt {
	return func(ps *PixelSort) {
		vPS := reflect.ValueOf(ps)

		vMethod := vPS.MethodByName(sortFunc)
		zero := reflect.Value{}
		if vMethod == zero {
			panic(fmt.Sprintf("sort func %s not found", sortFunc))
		}
		vPS.Elem().FieldByName("SorterFunc").Set(vMethod)

	}
}

func WithThresholdFuncString(thresholdFunc string) NewOpt {
	return func(ps *PixelSort) {
		vPS := reflect.ValueOf(ps)

		vMethod := vPS.MethodByName(thresholdFunc)
		zero := reflect.Value{}
		if vMethod == zero {
			panic(fmt.Sprintf("threshold func %s not found", thresholdFunc))
		}
		vPS.Elem().FieldByName("ThresholdFunc").Set(vMethod)

	}
}

func WithInvert(invert bool) NewOpt {
	return func(ps *PixelSort) {
		ps.invert = invert
	}
}
