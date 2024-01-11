package sqlite_test

import (
	"fmt"
	"testing"

	"github.com/senomas/gohtmx/stores"
	"github.com/stretchr/testify/assert"
)

func TestCRUDPrivilege(t *testing.T) {
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

	t.Run("find privilege", func(t *testing.T) {
		eprivileges := []stores.Privilege{
			{ID: 1, Name: rs("Admin"), Description: rs("Administrator")},
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, total, err := store.FindPrivileges(stores.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("populate privilege with dummy data", func(t *testing.T) {
		privileges := []stores.Privilege{}
		for i := 0; i < 100; i++ {
			privileges = append(privileges, stores.Privilege{
				Name:        rs(fmt.Sprintf("Demo-%d", i)),
				Description: rs(fmt.Sprintf("Demo %d", i)),
			})
		}
		actualPrivileges, err := store.AddPrivileges(privileges)
		assert.NoError(t, err)
		eprivileges := []stores.Privilege{}
		for i, p := range privileges {
			p.ID = 4 + int64(i)
			eprivileges = append(eprivileges, p)
		}
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.Equal(t, eprivileges, actualPrivileges)
	})

	t.Run("add non unique privilege", func(t *testing.T) {
		privileges := []stores.Privilege{
			{Name: rs("Root"), Description: rs("Im GROOT")},
			{Name: rs("User"), Description: rs("User")},
		}
		_, err := store.AddPrivileges(privileges)
		assert.ErrorContains(t, err, "duplicate record privilege.name 'User'")
	})

	t.Run("find privilege with offset and limit", func(t *testing.T) {
		eprivileges := []stores.Privilege{
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		for i := 0; i < 8; i++ {
			eprivileges = append(eprivileges, stores.Privilege{
				ID:          4 + int64(i),
				Name:        rs(fmt.Sprintf("Demo-%d", i)),
				Description: rs(fmt.Sprintf("Demo %d", i)),
			})
		}
		actualPrivileges, total, err := store.FindPrivileges(stores.PrivilegeFilter{}, 1, 10)
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
		err := store.DeletePrivileges(ids)
		assert.NoError(t, err)
	})

	t.Run("find privilege", func(t *testing.T) {
		eprivileges := []stores.Privilege{
			{ID: 1, Name: rs("Admin"), Description: rs("Administrator")},
			{ID: 2, Name: rs("User"), Description: rs("User")},
			{ID: 3, Name: rs("Guest"), Description: rs("Guest")},
		}
		actualPrivileges, total, err := store.FindPrivileges(stores.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(eprivileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, eprivileges, actualPrivileges)
	})
}

func rs(v string) *string {
	return &v
}
