package sqlite

import (
	"fmt"

	"github.com/senomas/gohtmx/store"
)

// FindUsers implements store.store.
func (s *SqliteAccountStore) FindUsers(f *store.UserFilter, offset int64, limit int) ([]*store.User, int64, error) {
	ctx := filter{}
	ctx.Int64("id", f.ID)
	ctx.String("name", f.Name)
	ctx.String("email", f.Email)

	if !s.ValidLimit(limit) {
		return nil, 0, fmt.Errorf("invalid limit %d", limit)
	}

	qry := "SELECT count(id) FROM user"
	qry = ctx.AppendWhere(qry)
	var total int64
	err := s.db.Get(&total, qry, ctx.args...)
	if err != nil {
		return nil, 0, err
	}
	users := []*store.User{}
	qry = "SELECT id, name, email, password FROM user"
	qry = ctx.AppendWhere(qry)
	qry += " LIMIT ? OFFSET ?"
	args := append(ctx.args, limit, offset)
	err = s.db.Select(&users, qry, args...)
	return users, total, err
}
