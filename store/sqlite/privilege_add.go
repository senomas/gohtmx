package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/gohtmx/store"
)

// AddPrivileges implements store.Store.
func (s *SqliteAccountStore) AddPrivileges(privileges []*store.Privilege) ([]*store.Privilege, error) {
	tx := s.db.MustBegin()
	defer tx.Rollback()
	ps, err := tx.PrepareNamed("INSERT INTO privilege (name, description) VALUES (:name, :description)")
	if err != nil {
		return nil, fmt.Errorf("error creating PrepareNamed: %v", err)
	}
	res := []*store.Privilege{}
	for _, privilege := range privileges {
		rs, err := ps.Exec(privilege)
		if err != nil {
			em := err.Error()
			if strings.HasPrefix(em, "UNIQUE constraint failed: ") {
				ks := em[26:]
				ka := strings.SplitN(ks, ".", 3)
				var v interface{}
				if len(ka) == 2 {
					switch ka[1] {
					case "name":
						v = *privilege.Name
					default:
						v = s.ValueString(privilege)
					}
				}
				return nil, fmt.Errorf("error insert privilege%s: duplicate record %s '%v'", s.ValueString(privilege), ks, v)
			}
			return nil, fmt.Errorf("error insert privilege%s: %v", s.ValueString(privilege), err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return res, fmt.Errorf("error insert privilege%s affected %v: %v", s.ValueString(privilege), affected, err)
		}
		if affected != 1 {
			return res, fmt.Errorf("error insert privilege%s affected %v", s.ValueString(privilege), affected)
		}
		id, err := rs.LastInsertId()
		if err != nil {
			return res, fmt.Errorf("error insert privilege%s get id: %v", s.ValueString(privilege), err)
		}
		privilege.ID = id
		res = append(res, privilege)
	}
	err = tx.Commit()
	return res, err
}
