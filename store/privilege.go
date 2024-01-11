package store

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

func (p *Privilege) SetName(v string) *Privilege {
	p.Name = &v
	return p
}

func (p *Privilege) SetDescription(v string) *Privilege {
	p.Description = &v
	return p
}
