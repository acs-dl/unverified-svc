package models

import (
	"github.com/acs-dl/unverified-svc/internal/data"
	"github.com/acs-dl/unverified-svc/resources"
	"strings"
)

func NewUserListResponse(users []data.User) UserListResponse {
	return UserListResponse{
		Data: newUserList(users),
	}
}

func NewUserResponse(user data.User) resources.UserResponse {
	return resources.UserResponse{
		Data: newUser(user),
	}
}

func newUserList(users []data.User) []resources.User {
	result := make([]resources.User, 0)

	for _, user := range users {
		result = append(result, newUser(user))
	}

	return result
}

func newUser(user data.User) resources.User {
	modules := make([]string, 0)
	for _, module := range strings.Split(user.Module, ",") {
		modules = append(modules, module)
	}

	return resources.User{
		Key: resources.NewKeyInt64(user.Id, resources.USER),
		Attributes: resources.UserAttributes{
			Module:    modules,
			ModuleId:  user.ModuleId,
			Submodule: strings.Split(user.Submodule, ","),
			CreatedAt: user.CreatedAt,
			Name:      user.Name,
			Username:  user.Username,
			Phone:     user.Phone,
			Email:     user.Email,
		},
	}
}

type UserListResponse struct {
	Meta  Meta             `json:"meta"`
	Data  []resources.User `json:"data"`
	Links *resources.Links `json:"links"`
}

type Meta struct {
	TotalCount int64 `json:"total_count"`
}
