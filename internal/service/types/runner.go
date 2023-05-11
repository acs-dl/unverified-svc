package types

import (
	"context"

	"github.com/acs-dl/unverified-svc/internal/config"
)

type Runner = func(context context.Context, config config.Config)
