package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Users interface {
	New() Users

	Upsert(user User) error
	Delete(user User) error
	Select() ([]User, error)
	Get() (*User, error)

	WithGroupedModulesAndSubmodules(modules *string) Users
	WithGroupedSubmodules(username, module *string) Users

	FilterByModules(modules ...string) Users
	FilterByUsernames(usernames ...string) Users
	FilterByPhones(phones ...string) Users
	FilterByEmails(emails ...string) Users

	SearchBy(search string) Users

	ResetFilters() Users

	Count() Users
	CountWithGroupedModules(modules *string) Users
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams, sortParams SortParams) Users
}

type User struct {
	Id        int64     `json:"-" db:"id" structs:"-"`
	Username  *string   `json:"username" db:"username" structs:"username,omitempty"`
	Phone     *string   `json:"phone" db:"phone" structs:"phone,omitempty"`
	Email     *string   `json:"email" db:"email" structs:"email,omitempty"`
	Name      *string   `json:"name" db:"name" structs:"name,omitempty"`
	Module    string    `json:"module" db:"module" structs:"module"`
	Submodule string    `json:"submodule" db:"submodule" structs:"submodule"`
	ModuleId  string    `json:"module_id" db:"module_id" structs:"module_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" structs:"-"`
}
