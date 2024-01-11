package sqlite_test

import (
	"fmt"
	"testing"

	"github.com/senomas/gohtmx/store"
	_ "github.com/senomas/gohtmx/store/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestCRUDPrivilege(t *testing.T) {
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

	t.Run("find privilege", func(t *testing.T) {
		eprivileges := []store.Privilege{
			{ID: 1, Name: rs("Admin"), Description: rs("Administrator")},
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, total, err := storeCtx.FindPrivileges(store.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("populate privilege with dummy data", func(t *testing.T) {
		privileges := []store.Privilege{}
		for i := 0; i < 100; i++ {
			privileges = append(privileges, store.Privilege{
				Name:        rs(fmt.Sprintf("Demo-%d", i)),
				Description: rs(fmt.Sprintf("Demo %d", i)),
			})
		}
		actualPrivileges, err := storeCtx.AddPrivileges(privileges)
		assert.NoError(t, err)
		eprivileges := []store.Privilege{}
		for i, p := range privileges {
			p.ID = 4 + int64(i)
			eprivileges = append(eprivileges, p)
		}
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("add non unique privilege", func(t *testing.T) {
		privileges := []store.Privilege{
			{Name: rs("Root"), Description: rs("Im GROOT")},
			{Name: rs("User"), Description: rs("User")},
		}
		_, err := storeCtx.AddPrivileges(privileges)
		assert.ErrorContains(t, err, "duplicate record privilege.name 'User'")
	})

	t.Run("find privilege with offset and limit", func(t *testing.T) {
		eprivileges := []store.Privilege{
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		for i := 0; i < 8; i++ {
			eprivileges = append(eprivileges, store.Privilege{
				ID:          4 + int64(i),
				Name:        rs(fmt.Sprintf("Demo-%d", i)),
				Description: rs(fmt.Sprintf("Demo %d", i)),
			})
		}
		actualPrivileges, total, err := storeCtx.FindPrivileges(store.PrivilegeFilter{}, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 103, total, "total")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("delete privilege dummy data", func(t *testing.T) {
		ids := []int64{}
		for i := 0; i < 100; i++ {
			ids = append(ids, 4+int64(i))
		}
		err := storeCtx.DeletePrivileges(ids)
		assert.NoError(t, err)
	})

	t.Run("find privilege", func(t *testing.T) {
		eprivileges := []store.Privilege{
			{ID: 1, Name: rs("Admin"), Description: rs("Administrator")},
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, total, err := storeCtx.FindPrivileges(store.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	/*
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
	*/
}

func rs(v string) *string {
	return &v
}
