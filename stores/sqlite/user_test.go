package sqlite_test

import (
	"testing"

	"github.com/senomas/gohtmx/stores"
	"github.com/stretchr/testify/assert"
)

func TestCRUDUser(t *testing.T) {
	store := stores.Get("sqlite")

	t.Run("populate privilege", func(t *testing.T) {
		privileges := []stores.Privilege{
			{Name: rs("Admin"), Description: rs("Administrator")},
			{Name: rs("User"), Description: rs("User")},
			{Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, err := store.AddPrivileges(privileges)
		assert.NoError(t, err)
		eprivileges := []stores.Privilege{}
		for i, p := range privileges {
			p.ID = 1 + int64(i)
			eprivileges = append(eprivileges, p)
		}
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("populate user", func(t *testing.T) {
		users := []stores.User{
			{
				Email:      rs("admin@cool.com"),
				Password:   rs("admin"),
				Name:       rs("Administrator"),
				Privileges: &[]stores.Privilege{{Name: rs("Admin")}},
			},
			{
				Email:      rs("user1@foo.com"),
				Password:   rs("user1"),
				Name:       rs("User 1"),
				Privileges: &[]stores.Privilege{{Name: rs("User")}},
			},
		}
		actualUsers, err := store.AddUsers(users)
		assert.NoError(t, err)
		eusers := []stores.User{}
		for i, u := range users {
			u.ID = 1 + int64(i)
			eprivileges := []stores.Privilege{}
			for _, p := range *u.Privileges {
				ep, err := store.GetPrivilegeByName(*p.Name)
				assert.NoError(t, err)
				eprivileges = append(eprivileges, *ep)
			}
			u.Privileges = &eprivileges
			eusers = append(eusers, u)
		}
		assert.Equal(t, len(eusers), len(actualUsers), "len")
		assert.Equal(t, eusers, actualUsers)
	})
}
