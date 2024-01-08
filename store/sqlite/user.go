package sqlite

import (
	"fmt"

	"github.com/senomas/gohtmx/store"
)

// GetUser implements store.StoreCtx.
func (s *SqliteStoreCtx) GetUser(id int64) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE id = ?", id)
	return &user, err
}

// FindUser implements store.StoreCtx.
func (s *SqliteStoreCtx) FindUser(f store.UserFilter, offset int64, limit int) ([]store.User, int64, error) {
	ctx := filterCtx{}
	ctx.Int64("id", f.ID)
	ctx.String("name", f.Name)

	qry := "SELECT count(id) FROM user"
	qry = ctx.AppendWhere(qry)
	fmt.Printf("USER FIND TOTAL [%s] : %+v\n", qry, ctx.args)
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
	fmt.Printf("USER FIND [%s] : %+v\n", qry, args)
	err = s.db.Select(&users, qry, args...)
	return users, total, err
}

// AddUser implements store.StoreCtx.
func (s *SqliteStoreCtx) AddUser(users []*store.User) error {
	tx := s.db.MustBegin()
	ps, err := tx.PrepareNamed("INSERT INTO user (name, email, password) VALUES (:name, :email, :password)")
	if err != nil {
		return fmt.Errorf("error creating PrepareNamed: %v", err)
	}
	for _, user := range users {
		rs, err := ps.Exec(user)
		if err != nil {
			return fmt.Errorf("error insert %+v: %v", user, err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return fmt.Errorf("error insert rows affected %v, %+v: %v", affected, user, err)
		}
		if affected != 1 {
			return fmt.Errorf("error insert rows affected %v, %+v", affected, user)
		}
		id, err := rs.LastInsertId()
		if err != nil {
			return fmt.Errorf("error insert rows get id %+v: %v", user, err)
		}
		user.ID = id
	}
	err = tx.Commit()
	return err
}

// DeleteUser implements store.StoreCtx.
func (*SqliteStoreCtx) DeleteUser(ids []int) error {
	panic("unimplemented")
}
