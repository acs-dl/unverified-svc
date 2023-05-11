package registrator

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/acs-dl/unverified-svc/internal/config"
	"github.com/acs-dl/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/running"
)

const serviceName = data.ModuleName + "-registrar"

type Registrar interface {
	Run(ctx context.Context)
	UnregisterModule() error
}

type registrar struct {
	logger *logan.Entry
	config config.RegistratorConfig
}

func NewRegistrar(cfg config.Config) Registrar {
	return &registrar{
		logger: cfg.Log().WithField("runner", serviceName),
		config: cfg.Registrator(),
	}
}

func (r *registrar) Run(ctx context.Context) {
	running.WithBackOff(
		ctx,
		r.logger,
		serviceName,
		r.registerModule,
		10*time.Minute,
		10*time.Minute,
		10*time.Minute,
	)
}
