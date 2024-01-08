package store

type Privilege struct {
	Name        *string
	Description *string
	ID          int64
}

type PrivilegeFilter struct {
	Name        FilterString
	Description FilterString
	ID          FilterInt64
}
