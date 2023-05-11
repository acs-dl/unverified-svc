package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.com/distributed_lab/acs/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetUsersRequest struct {
	pgdb.OffsetPageParams
	data.SortParams
	Search *string `filter:"search"`
	Module *string `filter:"module"`
}

func NewGetUsersRequest(r *http.Request) (*GetUsersRequest, error) {
	var request GetUsersRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	return &request, request.validate()
}

func (r *GetUsersRequest) validate() error {
	return validation.Errors{
		"page[sort]": validation.Validate(r.Param, validation.In("created_at", "name", "username")),
	}.Filter()
}
