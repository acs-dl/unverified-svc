package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetUserRequest struct {
	Module   *string `filter:"module"`
	Username *string `filter:"username"`
}

func NewGetUserRequest(r *http.Request) (*GetUserRequest, error) {
	var request GetUserRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	return &request, request.validate()
}

func (r *GetUserRequest) validate() error {
	return validation.Errors{
		"module":   validation.Validate(r.Module, validation.Required),
		"username": validation.Validate(r.Username, validation.Required),
	}.Filter()
}
