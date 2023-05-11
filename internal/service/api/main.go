package api

import (
	"context"
	"net/http"

	"github.com/acs-dl/unverified-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type api struct {
	cfg config.Config
}

func (r *api) run() error {
	router := r.apiRouter()

	if err := r.cfg.Copus().RegisterChi(router); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(r.cfg.Listener(), router)
}

func NewApiRouter(cfg config.Config) *api {
	return &api{cfg: cfg}
}

func Run(_ context.Context, cfg config.Config) {
	if err := NewApiRouter(cfg).run(); err != nil {
		panic(err)
	}
}
