package service

import (
	"context"

	"github.com/tortugait/decision-maker/internal/domain"
)

type (
	DecisionSrv interface {
		ShouldDoIt() bool
		GetFullDecision() (domain.Decision, error)
	}

	DeepSeekSrv interface {
		GenerateReason(
			ctx context.Context,
			statement string,
			ok bool,
		) (string, error)
	}
)
