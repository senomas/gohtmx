package sqlite

import (
	"github.com/senomas/gohtmx/store"
)

// GetPrivilege implements store.Store.
func (s *SqliteAccountStore) GetPrivilege(id int64) (*store.Privilege, error) {
	var privilege store.Privilege
	err := s.db.Get(&privilege, "SELECT id, name, description FROM privilege WHERE id = ?", id)
	return &privilege, err
}

// GetPrivilegeByName implements store.Store.
func (s *SqliteAccountStore) GetPrivilegeByName(name string) (*store.Privilege, error) {
	var privilege store.Privilege
	err := s.db.Get(&privilege, "SELECT id, name, description FROM privilege WHERE name = ?", name)
	return &privilege, err
}

// GetUserPrivileges implements store.Store.
func (s *SqliteAccountStore) GetUserPrivileges(userID int64) ([]store.UserPrivilege, error) {
	privileges := []store.UserPrivilege{}
	err := s.db.Select(&privileges, "SELECT p.id, p.name, p.description FROM privilege p JOIN user_privilege up ON p.id = up.privilege WHERE up.user = ?", userID)
	return privileges, err
}
