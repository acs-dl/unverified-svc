package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

type Users interface {
	New() Users

	Upsert(user User) error
	Delete(user User) error
	Select() ([]User, error)
	Get() (*User, error)

	FilterByModules(modules ...string) Users
	FilterByUsernames(usernames ...string) Users
	FilterByPhones(phones ...string) Users
	FilterByEmails(emails ...string) Users

	SearchBy(search string) Users
	GroupBy(column string) Users

	ResetFilters() Users

	Count() Users
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) Users
}

type User struct {
	Id        int64     `json:"-" db:"id" structs:"-"`
	Username  *string   `json:"username" db:"username" structs:"username,omitempty"`
	Phone     *string   `json:"phone" db:"phone" structs:"phone,omitempty"`
	Email     *string   `json:"email" db:"email" structs:"email,omitempty"`
	Name      *string   `json:"name" db:"name" structs:"name,omitempty"`
	Module    string    `json:"module" db:"module" structs:"module"`
	ModuleId  int64     `json:"module_id" db:"module_id" structs:"module_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" structs:"-"`
}
