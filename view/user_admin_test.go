package view

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/stores"
	_ "github.com/senomas/gohtmx/stores/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestUserAdmin(t *testing.T) {
	db_type := os.Getenv("DB_TYPE")
	if db_type == "" {
		db_type = "sqlite"
	}
	storeCtx := stores.Get(db_type)
	assert.NotNil(t, storeCtx)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ViewContext{}
			cc.Context = c
			cc.store = storeCtx
			return next(cc)
		}
	})
	e.GET("/user", UserAdminHandler)
	go func() {
		e.Start(":3000")
	}()
	defer e.Close()

	t.Run("initialize privilege", initPrivileges(storeCtx))
	t.Run("initialize user", initUsers(storeCtx))

	t.Run("GET /user", func(t *testing.T) {
		res, err := http.Get("http://localhost:3000/user?_o=10")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		estr := ``
		assert.Equal(t, estr, string(body))
	})
}

func initPrivileges(storeCtx stores.Store) func(t *testing.T) {
	rs := func(s string) *string { return &s }
	return func(t *testing.T) {
		newPrivileges := []stores.Privilege{
			{Name: rs("Admin"), Description: rs("Administrator")},
			{Name: rs("User"), Description: rs("User")},
			{Name: rs("Guest"), Description: rs("Guest")},
		}
		_, err := storeCtx.AddPrivileges(newPrivileges)
		assert.NoError(t, err)
	}
}

func initUsers(storeCtx stores.Store) func(t *testing.T) {
	rs := func(s string) *string { return &s }
	return func(t *testing.T) {
		newUsers := []stores.User{
			{Name: rs("Admin 1"), Email: rs("admin1@cool.com"), Password: stores.HashPassword("dodol123"), Privileges: &[]stores.Privilege{
				{Name: rs("Admin")},
				{Name: rs("User")},
			}},
			{Name: rs("User 1"), Email: rs("user1@foo.com"), Password: stores.HashPassword("dodol123"), Privileges: &[]stores.Privilege{
				{Name: rs("User")},
			}},
			{Name: rs("User 2"), Email: rs("user2@foo.com"), Password: stores.HashPassword("duren123"), Privileges: &[]stores.Privilege{
				{Name: rs("User")},
			}},
			{Name: rs("User 3"), Email: rs("user3@foo.com"), Password: stores.HashPassword("dodol123"), Privileges: &[]stores.Privilege{
				{Name: rs("User")},
			}},
		}
		_, err := storeCtx.AddUsers(newUsers)
		assert.NoError(t, err)
	}
}
