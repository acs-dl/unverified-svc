package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/unverified-svc/internal/service/api/models"
	"gitlab.com/distributed_lab/acs/unverified-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := UsersQ(r).WithGroupedSubmodules(request.Username, request.Module).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get unverified user from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		Log(r).WithError(err).Errorf("no such unverified user")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, models.NewUserResponse(*user))
}
