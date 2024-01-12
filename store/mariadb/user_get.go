package mariadb

import (
	"github.com/senomas/gohtmx/store"
)

// GetUser implements store.store.
func (s *MariadbAccountStore) GetUser(id int64) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	privileges := []*store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", id)
	user.Privileges = &privileges
	return &user, err
}

// GetUserByName implements store.store.
func (s *MariadbAccountStore) GetUserByName(name string) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	privileges := []*store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", user.ID)
	user.Privileges = &privileges
	return &user, err
}

// GetUserByEmail implements store.store.
func (s *MariadbAccountStore) GetUserByEmail(email string) (*store.User, error) {
	var user store.User
	err := s.db.Get(&user, "SELECT id, name, email, password FROM user WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	privileges := []*store.Privilege{}
	err = s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", user.ID)
	user.Privileges = &privileges
	return &user, err
}
