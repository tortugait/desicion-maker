package service

import (
	"math/rand"
	"time"

	"github.com/tortugait/decision-maker/internal/domain"
)

const maxOptions = 2

type (
	decisionSrv struct {
		deepSeekSrv DeepSeekSrv
	}
)

func NewDecisionSrv(deepSeekSrv DeepSeekSrv) decisionSrv {
	return decisionSrv{deepSeekSrv: deepSeekSrv}
}

func (s decisionSrv) ShouldDoIt() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src) //nolint:gosec

	return rng.Intn(maxOptions) == 1
}

func (s decisionSrv) GetFullDecision() (domain.Decision, error) {
	// TODO: implement me
	return domain.Decision{}, nil // fake return
}
