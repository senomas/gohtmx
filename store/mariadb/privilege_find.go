package mariadb

import (
	"fmt"

	"github.com/senomas/gohtmx/store"
)

// FindPrivileges implements store.Store.
func (s *MariadbAccountStore) FindPrivileges(
	f *store.PrivilegeFilter, offset int64, limit int,
) ([]*store.Privilege, int64, error) {
	ctx := filter{}
	ctx.Int64("id", f.ID)
	ctx.String("name", f.Name)
	ctx.String("description", f.Description)

	if !s.ValidLimit(limit) {
		return nil, 0, fmt.Errorf("invalid limit %d", limit)
	}

	qry := "SELECT count(id) FROM privilege"
	qry = ctx.AppendWhere(qry)
	var total int64
	err := s.db.Get(&total, qry, ctx.args...)
	if err != nil {
		return nil, 0, err
	}
	privileges := []*store.Privilege{}
	qry = "SELECT id, name, description FROM privilege"
	qry = ctx.AppendWhere(qry)
	qry += " LIMIT ? OFFSET ?"
	args := append(ctx.args, limit, offset)
	err = s.db.Select(&privileges, qry, args...)
	return privileges, total, err
}
