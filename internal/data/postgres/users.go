package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	usersTableName      = "users"
	usersModuleIdColumn = usersTableName + ".module_id"
	usersModuleColumn   = usersTableName + ".module"
)

type UsersQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

var selectedUsersTable = sq.Select("*").From(usersTableName)

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:            db.Clone(),
		selectBuilder: selectedUsersTable,
		deleteBuilder: sq.Delete(usersTableName),
	}
}

func (q UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q UsersQ) Upsert(user data.User) error {
	updateStmt, args := sq.Update(" ").
		Set("username", user.Username).
		Set("phone", user.Phone).
		Set("email", user.Email).
		Set("name", user.Name).
		Set("created_at", time.Now()).MustSql()

	query := sq.Insert(usersTableName).SetMap(structs.Map(user)).
		Suffix("ON CONFLICT (module_id, module, submodule) DO "+updateStmt, args...)
	fmt.Println(query.MustSql())
	return q.db.Exec(query)
}

func (q UsersQ) Delete() error {
	var deleted []data.User

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return errors.Errorf("no rows deleted")
	}

	return nil
}

func (q UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q UsersQ) FilterByModuleIds(moduleIds ...string) data.Users {
	equalModuleIds := sq.Eq{usersModuleIdColumn: moduleIds}
	q.selectBuilder = q.selectBuilder.Where(equalModuleIds)
	q.deleteBuilder = q.deleteBuilder.Where(equalModuleIds)

	return q
}

func (q UsersQ) FilterByModules(modules ...string) data.Users {
	equalModules := sq.Eq{usersModuleColumn: modules}
	q.selectBuilder = q.selectBuilder.Where(equalModules)
	q.deleteBuilder = q.deleteBuilder.Where(equalModules)

	return q
}

func (q UsersQ) SearchBy(search string) data.Users {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")

	q.selectBuilder = q.selectBuilder.Where(sq.Or{
		sq.ILike{"t.username": search},
		sq.ILike{"m.name": search},
	})

	return q
}

func (q UsersQ) WithGroupedModulesAndSubmodules(module *string) data.Users {
	selectGroupedUsers := sq.Select(
		"username",
		"MAX(created_at) as created_at",
		"string_agg(DISTINCT module, ',') as module",
		"string_agg(submodule, ',') as submodule").
		From(usersTableName).
		GroupBy("username")

	if module != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"module": *module})
	}

	q.selectBuilder = sq.Select("t.username, t.module, t.submodule, t.created_at, m.name, m.phone, m.email, m.id, m.module_id").
		FromSelect(selectGroupedUsers, "t").
		Join("(SELECT DISTINCT ON (username) username, name, phone, email, id, module_id FROM users) m ON m.username = t.username")

	return q
}

func (q UsersQ) WithGroupedSubmodules(username, module *string) data.Users {
	selectGroupedUsers := sq.Select(
		"username",
		"MAX(created_at) as created_at",
		"string_agg(submodule, ',') as submodule",
		"string_agg(DISTINCT module, ',') as module").
		From(usersTableName).
		GroupBy("username")

	if module != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"module": *module})
	}
	if username != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"username": *username})
	}

	q.selectBuilder = sq.Select("t.username, t.module, t.created_at, t.submodule, m.name, m.phone, m.email, m.id, m.module_id").
		FromSelect(selectGroupedUsers, "t").
		Join("(SELECT DISTINCT ON (username) username, name, phone, email, id, module_id FROM users) m ON m.username = t.username")

	return q
}

func (q UsersQ) CountWithGroupedModules(module *string) data.Users {
	selectGroupedUsers := sq.Select("username", "MAX(created_at) as created_at", "string_agg(DISTINCT module, ',') as module").
		From(usersTableName).
		GroupBy("username")

	if module != nil {
		selectGroupedUsers = selectGroupedUsers.Where(sq.Eq{"module": *module})
	}

	q.selectBuilder = sq.Select("COUNT (*)").
		FromSelect(selectGroupedUsers, "t").
		Join("(SELECT DISTINCT ON (username) username, name, phone, email, id, module_id FROM users) m ON m.username = t.username")

	return q
}

func (q UsersQ) GetTotalCount() (int64, error) {
	var count int64

	err := q.db.Get(&count, q.selectBuilder)

	return count, err
}

func (q UsersQ) Page(pageParams pgdb.OffsetPageParams, sortParams data.SortParams) data.Users {
	q.selectBuilder = pageParams.ApplyTo(q.selectBuilder, sortParams.Param)

	return q
}
