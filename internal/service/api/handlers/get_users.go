package handlers

import (
	"net/http"

	"github.com/acs-dl/unverified-svc/internal/data"
	"github.com/acs-dl/unverified-svc/internal/service/api/models"
	"github.com/acs-dl/unverified-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUsersRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	search := ""
	if request.Search != nil {
		search = *request.Search
	}

	totalCount, err := UsersQ(r).CountWithGroupedModules(request.Module).SearchBy(search).GetTotalCount()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to select to get total count from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	users, err := UsersQ(r).WithGroupedModulesAndSubmodules(request.Module).SearchBy(search).Page(request.OffsetPageParams, request.SortParams).Select()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to select unverified users from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewUserListResponse(users)
	response.Meta.TotalCount = totalCount
	response.Links = data.GetOffsetLinks(r, request.OffsetPageParams, request.SortParams)
	ape.Render(w, response)
}
