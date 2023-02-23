package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetUsersRequest struct {
	pgdb.OffsetPageParams

	Search *string `filter:"search"`
}

func NewGetUsersRequest(r *http.Request) (GetUsersRequest, error) {
	var request GetUsersRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
