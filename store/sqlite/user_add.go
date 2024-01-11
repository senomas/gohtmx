package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/gohtmx/store"
)

// AddUsers implements store.store.
func (s *SqliteAccountStore) AddUsers(users []store.User) ([]store.User, error) {
	tx := s.db.MustBegin()
	defer tx.Rollback()
	ps, err := tx.PrepareNamed("INSERT INTO user (name, email, password) VALUES (:name, :email, :password)")
	if err != nil {
		return nil, fmt.Errorf("error prepare insert into user: %v", err)
	}
	psp, err := tx.PrepareNamed("INSERT INTO user_privilege (user, privilege) VALUES (:user, :privilege)")
	if err != nil {
		return nil, fmt.Errorf("error prepare insert into user_privilege: %v", err)
	}
	res := []store.User{}
	for _, user := range users {
		rs, err := ps.Exec(user)
		if err != nil {
			em := err.Error()
			if strings.HasPrefix(em, "UNIQUE constraint failed: ") {
				ks := em[26:]
				ka := strings.SplitN(ks, ".", 3)
				var v interface{}
				if len(ka) == 2 {
					switch ka[1] {
					case "name":
						v = *user.Name
					case "email":
						v = *user.Email
					default:
						v = s.ValueString(user)
					}
				}
				return nil, fmt.Errorf("error insert user%s: duplicate record %s '%v'", s.ValueString(user), ks, v)
			}
			return nil, fmt.Errorf("error insert user%s: %v", s.ValueString(user), err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return res, fmt.Errorf("error insert user%s affected %v: %v", s.ValueString(user), affected, err)
		}
		if affected != 1 {
			return res, fmt.Errorf("error insert user%s affected %v", s.ValueString(user), affected)
		}
		id, err := rs.LastInsertId()
		if err != nil {
			return res, fmt.Errorf("error insert user%s get id: %v", s.ValueString(user), err)
		}
		user.ID = id
		if user.Privileges != nil {
			privileges := []store.Privilege{}
			type UserPrivilege struct {
				User      int64
				Privilege int64
			}
			for _, privilege := range *user.Privileges {
				err := tx.Get(&privilege, "SELECT id, name, description FROM privilege WHERE name = ?", privilege.Name)
				if err != nil {
					return res, fmt.Errorf("error get privilege name '%s': %v", *privilege.Name, err)
				}
				up := UserPrivilege{User: user.ID, Privilege: privilege.ID}
				rs, err := psp.Exec(up)
				if err != nil {
					return nil, fmt.Errorf("error insert user_privilege%s: %v", s.ValueString(up), err)
				}
				affected, err := rs.RowsAffected()
				if err != nil {
					return res, fmt.Errorf("error insert user_privilege%s affected %v: %v", s.ValueString(up), affected, err)
				}
				if affected != 1 {
					return res, fmt.Errorf("error insert user_privilege%s affected %v", s.ValueString(up), affected)
				}
				privileges = append(privileges, privilege)
			}
			user.Privileges = &privileges
		}
		res = append(res, user)
	}
	err = tx.Commit()
	return res, err
}
