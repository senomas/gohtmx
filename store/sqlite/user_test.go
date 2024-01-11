package sqlite_test

import (
	"testing"

	"github.com/senomas/gohtmx/store"
	_ "github.com/senomas/gohtmx/store/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestCRUDUser(t *testing.T) {
	storeCtx := store.Get("sqlite")

	t.Run("populate privilege", func(t *testing.T) {
		privileges := []store.Privilege{
			{Name: rs("Admin"), Description: rs("Administrator")},
			{Name: rs("User"), Description: rs("User")},
			{Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, err := storeCtx.AddPrivileges(privileges)
		assert.NoError(t, err)
		eprivileges := []store.Privilege{}
		for i, p := range privileges {
			p.ID = 1 + int64(i)
			eprivileges = append(eprivileges, p)
		}
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("populate user", func(t *testing.T) {
		users := []store.User{
			{
				Email:      rs("admin@cool.com"),
				Password:   rs("admin"),
				Name:       rs("Administrator"),
				Privileges: &[]store.Privilege{{Name: rs("Admin")}},
			},
			{
				Email:      rs("user1@foo.com"),
				Password:   rs("user1"),
				Name:       rs("User 1"),
				Privileges: &[]store.Privilege{{Name: rs("User")}},
			},
		}
		actualUsers, err := storeCtx.AddUsers(users)
		assert.NoError(t, err)
		eusers := []store.User{}
		for i, u := range users {
			u.ID = 1 + int64(i)
			eprivileges := []store.Privilege{}
			for _, p := range *u.Privileges {
				ep, err := storeCtx.GetPrivilegeByName(*p.Name)
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
