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
		sq.ILike{"t.username": search},
		sq.ILike{"m.name": search},
		//sq.ILike{"m.phone": search},
		//sq.ILike{"m.email": search},
	})

	return q
}

func (q *UsersQ) WithGroupedModules(modules *string) data.Users {
	selectGroupedUsers := sq.Select("username", "MAX(created_at) as created_at", "string_agg(module, ',') as module").
		From(usersTableName).
		GroupBy("username")

	if modules != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"module": *modules})
	}

	q.sql = sq.Select("t.username, t.module, t.created_at, m.name, m.phone, m.email, m.id, m.module_id").
		FromSelect(selectGroupedUsers, "t").
		Join("(SELECT DISTINCT ON (username) username, name, phone, email, id, module_id FROM users) m ON m.username = t.username")

	return q
}

func (q *UsersQ) Count() data.Users {
	q.sql = sq.Select("COUNT (*)").From(usersTableName)

	return q
}

func (q *UsersQ) CountWithGroupedModules(module *string) data.Users {
	selectGroupedUsers := sq.Select("username", "MAX(created_at) as created_at", "string_agg(module, ',') as module").
		From(usersTableName).
		GroupBy("username")

	if module != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"module": *module})
	}

	q.sql = sq.Select("COUNT (*)").
		FromSelect(selectGroupedUsers, "t").
		Join("(SELECT DISTINCT ON (username) username, name, phone, email, id, module_id FROM users) m ON m.username = t.username")

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

func (q *UsersQ) Page(pageParams pgdb.OffsetPageParams, sortParams data.SortParams) data.Users {
	q.sql = pageParams.ApplyTo(q.sql, sortParams.Param)

	return q
}
