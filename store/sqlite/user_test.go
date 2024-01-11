package sqlite_test

import (
	"testing"

	"github.com/senomas/gohtmx/store"
	"github.com/stretchr/testify/assert"
)

func TestCRUDUser(t *testing.T) {
	accountStore := store.GetAccountStore("sqlite")

	t.Run("populate privilege", func(t *testing.T) {
		privileges := []*store.Privilege{
			(&store.Privilege{}).SetName("Admin").SetDescription("Administrator"),
			(&store.Privilege{}).SetName("User").SetDescription("User"),
			(&store.Privilege{}).SetName("Guest").SetDescription("Guest"),
		}
		actualPrivileges, err := accountStore.AddPrivileges(privileges)
		assert.NoError(t, err)
		for i, p := range privileges {
			p.ID = 1 + int64(i)
		}
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.Equal(t, privileges, actualPrivileges)
	})

	t.Run("populate user", func(t *testing.T) {
		users := []*store.User{
			(&store.User{}).
				SetName("Administrator").
				SetEmail("admin@cool.com").
				SetPassword("admin").
				AddPrivilege((&store.Privilege{}).SetName("Admin")),
			(&store.User{}).
				SetName("User 1").
				SetEmail("user1@foo.com").
				SetPassword("user1").
				AddPrivilege((&store.Privilege{}).SetName("User")),
		}
		actualUsers, err := accountStore.AddUsers(users)
		assert.NoError(t, err)
		for i, u := range users {
			u.ID = 1 + int64(i)
			eprivileges := []*store.Privilege{}
			for _, p := range *u.Privileges {
				ep, err := accountStore.GetPrivilegeByName(*p.Name)
				assert.NoError(t, err)
				eprivileges = append(eprivileges, ep)
			}
			u.Privileges = &eprivileges
		}
		assert.Equal(t, len(users), len(actualUsers), "len")
		assert.Equal(t, users, actualUsers)
	})
}
