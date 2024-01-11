package store

import "fmt"

type AccountStore interface {
	Close() error

	GetUser(id int64) (*User, error)
	GetUserByName(name string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	FindUsers(*UserFilter, int64, int) ([]*User, int64, error)
	AddUsers(users []*User) ([]*User, error)
	UpdateUser(user *User) error
	DeleteUsers(ids []int64) error

	GetPrivilege(id int64) (*Privilege, error)
	GetPrivilegeByName(name string) (*Privilege, error)
	FindPrivileges(*PrivilegeFilter, int64, int) ([]*Privilege, int64, error)
	AddPrivileges(privileges []*Privilege) ([]*Privilege, error)
	DeletePrivileges(ids []int64) error

	GetUserPrivileges(userID int64) ([]UserPrivilege, error)
}

var accountStores = map[string]func() AccountStore{}

func AddAccountStore(name string, store func() AccountStore) {
	accountStores[name] = store
}

func GetAccountStore(name string) AccountStore {
	v := accountStores[name]
	if v == nil {
		panic(fmt.Sprintf("No implementation for '%s'", name))
	}
	return v()
}
