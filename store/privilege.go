package store

type Privilege struct {
	Name        *string
	Description *string
	ID          *int64
}

type PrivilegeFilter struct {
	Name        FilterString
	Description FilterString
	ID          FilterInt64
}

func (p *Privilege) SetID(v int64) *Privilege {
	p.ID = &v
	return p
}

func (p *Privilege) SetName(v string) *Privilege {
	p.Name = &v
	return p
}

func (p *Privilege) SetDescription(v string) *Privilege {
	p.Description = &v
	return p
}
