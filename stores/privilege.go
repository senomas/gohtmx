package stores

type Privilege struct {
	Name        *string
	Description *string
	ID          int64
}

type UserPrivilege struct {
	Name        *string
	Description *string
	UserID      int64
	ID          int64
}

type PrivilegeFilter struct {
	Name        FilterString
	Description FilterString
	ID          FilterInt64
}