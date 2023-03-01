package models

import (
	"gitlab.com/distributed_lab/acs/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/acs/unverified-svc/resources"
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
	return resources.User{
		Key: resources.NewKeyInt64(user.Id, resources.USER),
		Attributes: resources.UserAttributes{
			Module:    user.Module,
			ModuleId:  user.ModuleId,
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
