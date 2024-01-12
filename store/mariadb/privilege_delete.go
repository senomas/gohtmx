package mariadb

import (
	"fmt"
)

// DeletePrivileges implements store.Store.
func (s *MariadbAccountStore) DeletePrivileges(ids []int64) error {
	tx := s.db.MustBegin()
	defer tx.Rollback()
	qry := "DELETE FROM privilege WHERE id IN ("
	args := []interface{}{}
	for i, id := range ids {
		if i > 0 {
			qry += ","
		}
		qry += "?"
		args = append(args, id)
	}
	qry += ")"
	rs, err := tx.Exec(qry, args...)
	if err != nil {
		em := err.Error()
		if em == "FOREIGN KEY constraint failed" {
			return fmt.Errorf("error delete privilege.id%s: record in use", s.ValueString(ids))
		}
		return fmt.Errorf("error delete privilege.id%s: %v", s.ValueString(ids), err)
	}
	affected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("error delete privilege.id%s affected: %v", s.ValueString(ids), err)
	}
	if affected != int64(len(ids)) {
		return fmt.Errorf("error delete privilege.id%s affected %v", s.ValueString(ids), affected)
	}
	err = tx.Commit()
	return err
}
