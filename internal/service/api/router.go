package api

import (
	"fmt"

	auth "github.com/acs-dl/auth-svc/middlewares"
	"github.com/acs-dl/unverified-svc/internal/data"
	"github.com/acs-dl/unverified-svc/internal/data/postgres"
	"github.com/acs-dl/unverified-svc/internal/service/api/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (r *api) apiRouter() chi.Router {
	router := chi.NewRouter()

	logger := r.cfg.Log().WithField("service", fmt.Sprintf("%s-api", data.ModuleName))

	secret := r.cfg.JwtParams().Secret

	router.Use(
		ape.RecoverMiddleware(logger),
		ape.LoganMiddleware(logger),
		ape.CtxMiddleware(
			handlers.CtxLog(logger),
			handlers.CtxUsersQ(postgres.NewUsersQ(r.cfg.DB())),
		),
	)

	router.Route("/integrations/unverified-svc", func(r chi.Router) {
		// configure endpoints here
		r.Get("/user_roles", handlers.GetUserRolesMap) // comes from orchestrator

		r.Route("/users", func(r chi.Router) {
			r.With(auth.Jwt(secret, data.ModuleName, []string{"read", "write"}...)).
				Get("/", handlers.GetUsers)
		})
		r.Route("/user", func(r chi.Router) {
			r.With(auth.Jwt(secret, data.ModuleName, []string{"read", "write"}...)).
				Get("/", handlers.GetUser)
		})
	})

	return router
}
