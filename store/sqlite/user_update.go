package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/gohtmx/store"
)

func (s *SqliteAccountStore) UpdateUser(user *store.User) error {
	updates := []string{}
	args := []interface{}{}
	if user.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *user.Name)
	}
	if user.Email != nil {
		updates = append(updates, "email = ?")
		args = append(args, *user.Email)
	}
	if user.Password != nil {
		updates = append(updates, "password = ?")
		args = append(args, store.HashPassword(*user.Password))
	}
	tx := s.db.MustBegin()
	defer tx.Rollback()
	if len(updates) > 0 {
		qry := "UPDATE user SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		args = append(args, user.ID)
		rs, err := tx.Exec(qry, args...)
		if err != nil {
			return fmt.Errorf("error update user %s: %v", qry, err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return fmt.Errorf("error update user%s affected: %v", s.ValueString(user), err)
		}
		if affected != 1 {
			return fmt.Errorf("error update user%s affected %v", s.ValueString(user), affected)
		}
	}
	if user.Privileges != nil {
		npname := []interface{}{}
		npid := []int64{}
		qry := "SELECT id FROM privilege WHERE name IN ("
		for i, privilege := range *user.Privileges {
			if i > 0 {
				qry += ","
			}
			qry += "?"
			npname = append(npname, *privilege.Name)
		}
		qry += ")"
		err := tx.Select(&npid, qry, npname...)
		if err != nil {
			return fmt.Errorf("error select privilege '%s' %+v: %v", qry, npname, err)
		}
		opid := []int64{}
		qry = "SELECT privilege FROM user_privilege WHERE user = ?"
		err = tx.Select(&opid, qry, user.ID)
		if err != nil {
			return fmt.Errorf("error select user_privilege '%s' %+v: %v", qry, user.ID, err)
		}
		ipid := []int64{}
		rpid := []int64{}
		for _, n := range npid {
			found := false
			for _, o := range opid {
				if n == o {
					found = true
				}
			}
			if !found {
				ipid = append(ipid, n)
			}
		}
		for _, o := range opid {
			found := false
			for _, n := range npid {
				if n == o {
					found = true
				}
			}
			if !found {
				rpid = append(rpid, o)
			}
		}
		type UserPrivilege struct {
			User      int64
			Privilege int64
		}
		if len(ipid) > 0 {
			ps, err := tx.PrepareNamed("INSERT INTO user_privilege (user, privilege) VALUES (:user, :privilege)")
			if err != nil {
				return fmt.Errorf("error prepare insert into user_privilege: %v", err)
			}
			for _, n := range ipid {
				up := UserPrivilege{User: *user.ID, Privilege: n}
				rs, err := ps.Exec(up)
				if err != nil {
					return fmt.Errorf("error insert user_privilege%s: %v", s.ValueString(up), err)
				}
				affected, err := rs.RowsAffected()
				if err != nil {
					return fmt.Errorf("error insert user_privilege%s affected %v: %v", s.ValueString(up), affected, err)
				}
				if affected != 1 {
					return fmt.Errorf("error insert user_privilege%s affected %v", s.ValueString(up), affected)
				}
			}
		}
		if len(rpid) > 0 {
			ps, err := tx.PrepareNamed("DELETE FROM user_privilege WHERE user = :user AND privilege = :privilege")
			if err != nil {
				return fmt.Errorf("error prepare delete user_privilege: %v", err)
			}
			for _, r := range rpid {
				up := UserPrivilege{User: *user.ID, Privilege: r}
				rs, err := ps.Exec(up)
				if err != nil {
					return fmt.Errorf("error delete user_privilege%s: %v", s.ValueString(up), err)
				}
				affected, err := rs.RowsAffected()
				if err != nil {
					return fmt.Errorf("error delete user_privilege%s affected %v: %v", s.ValueString(up), affected, err)
				}
				if affected != 1 {
					return fmt.Errorf("error delete user_privilege%s affected %v", s.ValueString(up), affected)
				}
			}
		}
	}
	err := tx.Commit()
	return err
}
