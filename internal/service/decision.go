package service

import (
	"math/rand"
	"time"
)

func ShouldDoIt() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src) //nolint:gosec

	return rng.Intn(2) == 1 //nolint:gomnd
}
