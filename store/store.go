package store

type StoreCtx interface {
	GetUser(id int64) (*User, error)
	FindUser(UserFilter, int64, int) ([]User, int64, error)
	AddUser(users []*User) error
	DeleteUser(ids []int) error
}
