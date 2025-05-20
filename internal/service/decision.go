package service

import (
	"math/rand"
	"time"
)

const maxOptions = 2

func ShouldDoIt() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src) //nolint:gosec

	return rng.Intn(maxOptions) == 1
}
