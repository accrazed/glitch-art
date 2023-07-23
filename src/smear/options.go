package smear

import "math/rand"

type SmearOpt func(*Smearer)

func WithSeed(seed int64) SmearOpt {
	return func(s *Smearer) {
		s.r = rand.New(rand.NewSource(seed))
	}
}

func WithStrength(strength uint) SmearOpt {
	return func(s *Smearer) {
		s.strength = strength
	}
}

func WithSmearPos(smearPos int) SmearOpt {
	return func(s *Smearer) {
		s.smearPos = smearPos
	}
}

func WithSmearLen(smearLen int) SmearOpt {
	return func(s *Smearer) {
		s.smearLen = smearLen
	}
}
