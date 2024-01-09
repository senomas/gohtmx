package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/gohtmx/store"
)

// GetUser implements store.StoreCtx.
func (s *SqliteStoreCtx) GetUser(id int64) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	privileges := []store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", id)
	user.Privileges = &privileges
	return &user, err
}

// GetUserByName implements store.StoreCtx.
func (s *SqliteStoreCtx) GetUserByName(name string) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	privileges := []store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", user.ID)
	user.Privileges = &privileges
	return &user, err
}

// GetUserByEmail implements store.StoreCtx.
func (s *SqliteStoreCtx) GetUserByEmail(email string) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	privileges := []store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", user.ID)
	user.Privileges = &privileges
	return &user, err
}

// FindUsers implements store.StoreCtx.
func (s *SqliteStoreCtx) FindUsers(f store.UserFilter, offset int64, limit int) ([]store.User, int64, error) {
	ctx := filterCtx{}
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
	users := []store.User{}
	qry = "SELECT id, name, email, password FROM user"
	qry = ctx.AppendWhere(qry)
	qry += " LIMIT ? OFFSET ?"
	args := append(ctx.args, limit, offset)
	err = s.db.Select(&users, qry, args...)
	return users, total, err
}

// AddUsers implements store.StoreCtx.
func (s *SqliteStoreCtx) AddUsers(users []store.User) ([]store.User, error) {
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

func (s *SqliteStoreCtx) UpdateUser(user store.User) error {
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
				up := UserPrivilege{User: user.ID, Privilege: n}
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
				up := UserPrivilege{User: user.ID, Privilege: r}
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

// DeleteUsers implements store.StoreCtx.
func (s *SqliteStoreCtx) DeleteUsers(ids []int64) error {
	tx := s.db.MustBegin()
	defer tx.Rollback()
	qry := "DELETE FROM user WHERE id IN ("
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
			return fmt.Errorf("error delete user.id%s: record in use", s.ValueString(ids))
		}
		return fmt.Errorf("error delete user.id%s: %v", s.ValueString(ids), err)
	}
	affected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("error delete user.id%s affected: %v", s.ValueString(ids), err)
	}
	if affected != int64(len(ids)) {
		return fmt.Errorf("error delete user.id%s affected %v", s.ValueString(ids), affected)
	}
	err = tx.Commit()
	return err
}
