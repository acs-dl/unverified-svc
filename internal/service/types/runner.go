package types

import (
	"context"

	"gitlab.com/distributed_lab/acs/unverified-svc/internal/config"
)

type Runner = func(context context.Context, config config.Config)
