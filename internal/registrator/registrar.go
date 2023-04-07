package registrator

import (
	"context"

	"gitlab.com/distributed_lab/acs/unverified-svc/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewRegistrar(cfg).Run(ctx)
}
