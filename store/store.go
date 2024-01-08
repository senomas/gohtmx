package store

type StoreCtx interface {
	Close() error

	GetUser(id int64) (*User, error)
	GetUserByName(name string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	FindUsers(UserFilter, int64, int) ([]User, int64, error)
	AddUsers(users []User) ([]User, error)
	UpdateUser(user User) error
	DeleteUsers(ids []int64) error

	GetPrivilege(id int64) (*Privilege, error)
	FindPrivileges(PrivilegeFilter, int64, int) ([]Privilege, int64, error)
	AddPrivileges(privileges []Privilege) ([]Privilege, error)
	DeletePrivileges(ids []int) error
}
