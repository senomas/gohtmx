package sqlite

import (
	"fmt"

	"github.com/senomas/gohtmx/store"
)

// GetPrivilege implements store.StoreCtx.
func (s *SqliteStoreCtx) GetPrivilege(id int64) (*store.Privilege, error) {
	var privilege store.Privilege
	err := s.db.Get(&privilege, "SELECT id, name, description FROM privilege WHERE id = ?", id)
	return &privilege, err
}

// FindPrivileges implements store.StoreCtx.
func (s *SqliteStoreCtx) FindPrivileges(f store.PrivilegeFilter, offset int64, limit int) ([]store.Privilege, int64, error) {
	ctx := filterCtx{}
	ctx.Int64("id", f.ID)
	ctx.String("name", f.Name)
	ctx.String("description", f.Description)

	qry := "SELECT count(id) FROM privilege"
	qry = ctx.AppendWhere(qry)
	var total int64
	err := s.db.Get(&total, qry, ctx.args...)
	if err != nil {
		return nil, 0, err
	}
	privileges := []store.Privilege{}
	qry = "SELECT id, name, description FROM privilege"
	qry = ctx.AppendWhere(qry)
	qry += " LIMIT ? OFFSET ?"
	args := append(ctx.args, limit, offset)
	err = s.db.Select(&privileges, qry, args...)
	return privileges, total, err
}

// AddPrivileges implements store.StoreCtx.
func (s *SqliteStoreCtx) AddPrivileges(privileges []store.Privilege) ([]store.Privilege, error) {
	tx := s.db.MustBegin()
	ps, err := tx.PrepareNamed("INSERT INTO privilege (name, description) VALUES (:name, :description)")
	if err != nil {
		return nil, fmt.Errorf("error creating PrepareNamed: %v", err)
	}
	res := []store.Privilege{}
	for _, privilege := range privileges {
		rs, err := ps.Exec(privilege)
		if err != nil {
			return res, fmt.Errorf("error insert %+v: %v", privilege, err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return res, fmt.Errorf("error insert rows affected %v, %+v: %v", affected, privilege, err)
		}
		if affected != 1 {
			return res, fmt.Errorf("error insert rows affected %v, %+v", affected, privilege)
		}
		id, err := rs.LastInsertId()
		if err != nil {
			return res, fmt.Errorf("error insert rows get id %+v: %v", privilege, err)
		}
		privilege.ID = id
		res = append(res, privilege)
	}
	err = tx.Commit()
	return res, err
}

// DeletePrivileges implements store.StoreCtx.
func (*SqliteStoreCtx) DeletePrivileges(ids []int) error {
	panic("unimplemented")
}
