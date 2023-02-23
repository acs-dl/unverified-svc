package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"strings"
	"time"
)

const usersTableName = "users"

type UsersQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedUsersTable = sq.Select("*").From(usersTableName)

var usersColumns = []string{
	usersTableName + ".id",
	usersTableName + ".username",
	usersTableName + ".phone",
	usersTableName + ".email",
	usersTableName + ".module",
	usersTableName + ".name",
	usersTableName + ".created_at",
}

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:  db.Clone(),
		sql: selectedUsersTable,
	}
}

func (q *UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q *UsersQ) Upsert(user data.User) error {
	updateStmt, args := sq.Update(" ").
		Set("created_at", time.Now()).MustSql()

	query := sq.Insert(usersTableName).SetMap(structs.Map(user)).
		Suffix("ON CONFLICT (module_id, module) DO "+updateStmt, args...)

	return q.db.Exec(query)
}

func (q *UsersQ) Delete(user data.User) error {
	query := sq.Delete(usersTableName).Where(
		sq.Eq{
			usersTableName + ".module":    user.Module,
			usersTableName + ".module_id": user.ModuleId,
		},
	)

	result, err := q.db.ExecWithResult(query)
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.Errorf("no users for module `%s` with module id `%d`", user.Module, user.ModuleId)
	}

	return nil
}

func (q *UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *UsersQ) FilterByModules(modules ...string) data.Users {
	q.sql = q.sql.Where(sq.Eq{usersTableName + ".module": modules})

	return q
}

func (q *UsersQ) FilterByUsernames(usernames ...string) data.Users {
	q.sql = q.sql.Where(sq.Eq{usersTableName + ".username": usernames})

	return q
}

func (q *UsersQ) FilterByPhones(phones ...string) data.Users {
	q.sql = q.sql.Where(sq.Eq{usersTableName + ".phone": phones})

	return q
}

func (q *UsersQ) FilterByEmails(emails ...string) data.Users {
	q.sql = q.sql.Where(sq.Eq{usersTableName + ".email": emails})

	return q
}

func (q *UsersQ) SearchBy(search string) data.Users {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")

	q.sql = q.sql.Where(sq.Or{
		sq.ILike{usersTableName + ".username": search},
		sq.ILike{usersTableName + ".phone": search},
		sq.ILike{usersTableName + ".email": search},
		sq.ILike{usersTableName + ".name": search},
	})

	return q
}

func (q *UsersQ) GroupBy(column string) data.Users {
	q.sql = q.sql.GroupBy(column)

	return q
}

func (q *UsersQ) Count() data.Users {
	q.sql = sq.Select("COUNT (*)").From(usersTableName)

	return q
}

func (q *UsersQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.sql)

	return count, err
}

func (q *UsersQ) ResetFilters() data.Users {
	q.sql = selectedUsersTable

	return q
}

func (q *UsersQ) Page(pageParams pgdb.OffsetPageParams) data.Users {
	q.sql = pageParams.ApplyTo(q.sql, "username")

	return q
}
