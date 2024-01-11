package sqlite_test

import (
	"fmt"
	"testing"

	"github.com/senomas/gohtmx/store"
	"github.com/stretchr/testify/assert"
)

func TestCRUDPrivilege(t *testing.T) {
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
			p.SetID(1 + int64(i))
		}
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.Equal(t, privileges, actualPrivileges)
	})

	t.Run("find privilege", func(t *testing.T) {
		privileges := []*store.Privilege{
			(&store.Privilege{}).SetID(1).SetName("Admin").SetDescription("Administrator"),
			(&store.Privilege{}).SetID(2).SetName("User").SetDescription("User"),
			(&store.Privilege{}).SetID(3).SetName("Guest").SetDescription("Guest"),
		}
		actualPrivileges, total, err := accountStore.FindPrivileges(&store.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, privileges, actualPrivileges)
	})

	t.Run("populate privilege with dummy data", func(t *testing.T) {
		privileges := []*store.Privilege{}
		for i := 0; i < 100; i++ {
			privileges = append(privileges,
				(&store.Privilege{}).
					SetName(fmt.Sprintf("Demo-%d", i)).SetDescription(fmt.Sprintf("Demo %d", i)))
		}
		actualPrivileges, err := accountStore.AddPrivileges(privileges)
		assert.NoError(t, err)
		for i, p := range privileges {
			p.SetID(4 + int64(i))
		}
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.Equal(t, privileges, actualPrivileges)
	})

	t.Run("add non unique privilege", func(t *testing.T) {
		privileges := []*store.Privilege{
			(&store.Privilege{}).SetName("Root").SetDescription("Im GROOT"),
			(&store.Privilege{}).SetName("User").SetDescription("User"),
		}
		_, err := accountStore.AddPrivileges(privileges)
		assert.ErrorContains(t, err, "duplicate record privilege.name 'User'")
	})

	t.Run("find privilege with offset and limit", func(t *testing.T) {
		privileges := []*store.Privilege{
			(&store.Privilege{}).SetID(2).SetName("User").SetDescription("User"),
			(&store.Privilege{}).SetID(3).SetName("Guest").SetDescription("Guest"),
		}
		for i := 0; i < 8; i++ {
			privileges = append(privileges,
				(&store.Privilege{}).
					SetID(4+int64(i)).
					SetName(fmt.Sprintf("Demo-%d", i)).
					SetDescription(fmt.Sprintf("Demo %d", i)))
		}
		actualPrivileges, total, err := accountStore.FindPrivileges(&store.PrivilegeFilter{}, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 103, total, "total")
		assert.Equal(t, privileges, actualPrivileges)
	})

	t.Run("delete privilege dummy data", func(t *testing.T) {
		ids := []int64{}
		for i := 0; i < 100; i++ {
			ids = append(ids, 4+int64(i))
		}
		err := accountStore.DeletePrivileges(ids)
		assert.NoError(t, err)
	})

	t.Run("find privilege", func(t *testing.T) {
		privileges := []*store.Privilege{
			(&store.Privilege{}).SetID(1).SetName("Admin").SetDescription("Administrator"),
			(&store.Privilege{}).SetID(2).SetName("User").SetDescription("User"),
			(&store.Privilege{}).SetID(3).SetName("Guest").SetDescription("Guest"),
		}
		actualPrivileges, total, err := accountStore.FindPrivileges(&store.PrivilegeFilter{}, 0, 100)
		assert.NoError(t, err)
		assert.Equal(t, len(privileges), len(actualPrivileges), "len")
		assert.EqualValues(t, 3, total, "total")
		assert.Equal(t, privileges, actualPrivileges)
	})
}
