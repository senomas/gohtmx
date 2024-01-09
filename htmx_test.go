package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
	"github.com/stretchr/testify/assert"
)

func TestTemplates(t *testing.T) {
	db_type := os.Getenv("DB_TYPE")
	if db_type == "" {
		db_type = "sqlite"
	}
	storeCtx := store.Get(db_type)
	assert.NotNil(t, storeCtx)
	e := echo.New()
	setupRouter(storeCtx, e)
	go func() {
		e.Start(":1323")
	}()
	defer e.Close()

	t.Run("initialize privilege", initPrivileges(storeCtx))
	t.Run("initialize user", initUsers(storeCtx))

	t.Run("GET /user/1", func(t *testing.T) {
		res, err := http.Get("http://localhost:1323/user/1")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body := map[string]interface{}{}
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		estr := MustSerialize(t, map[string]interface{}{
			"Email": "admin1@cool.com",
			"ID":    1,
			"Name":  "Admin 1",
			"Privileges": []map[string]interface{}{
				{
					"Description": "Administrator",
					"ID":          1,
					"Name":        "Admin",
				},
				{
					"Description": "User",
					"ID":          2,
					"Name":        "User",
				},
			},
		})

		assert.Equal(t, estr, MustSerialize(t, StripRow(t, body)))
	})

	t.Run("GET /user/name/Admin 1", func(t *testing.T) {
		res, err := http.Get("http://localhost:1323/user/name/Admin%201")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body := map[string]interface{}{}
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		estr := MustSerialize(t, map[string]interface{}{
			"Email": "admin1@cool.com",
			"ID":    1,
			"Name":  "Admin 1",
			"Privileges": []map[string]interface{}{
				{
					"Description": "Administrator",
					"ID":          1,
					"Name":        "Admin",
				},
				{
					"Description": "User",
					"ID":          2,
					"Name":        "User",
				},
			},
		})

		assert.Equal(t, estr, MustSerialize(t, StripRow(t, body)))
	})

	t.Run("GET /user/email/admin1@cool.com", func(t *testing.T) {
		res, err := http.Get("http://localhost:1323/user/email/admin1@cool.com")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body := map[string]interface{}{}
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		estr := MustSerialize(t, map[string]interface{}{
			"Email": "admin1@cool.com",
			"ID":    1,
			"Name":  "Admin 1",
			"Privileges": []map[string]interface{}{
				{
					"Description": "Administrator",
					"ID":          1,
					"Name":        "Admin",
				},
				{
					"Description": "User",
					"ID":          2,
					"Name":        "User",
				},
			},
		})

		assert.Equal(t, estr, MustSerialize(t, StripRow(t, body)))
	})

	t.Run("POST /user demo", func(t *testing.T) {
		rs := func(s string) *string { return &s }
		newUsers := []store.User{
			{Name: rs("Demo"), Email: rs("demo@foo.com"), Password: store.HashPassword("dodol123"), Privileges: &[]store.Privilege{
				{Name: rs("User")},
			}},
		}
		s, _ := json.Marshal(newUsers)
		res, err := http.Post("http://localhost:1323/user", "application/json", bytes.NewReader(s))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body := []map[string]interface{}{}
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		estr := MustSerialize(t, []map[string]interface{}{
			{
				"Email": "demo@foo.com",
				"ID":    5,
				"Name":  "Demo",
				"Privileges": []map[string]interface{}{
					{
						"Description": "User",
						"ID":          2,
						"Name":        "User",
					},
				},
			},
		})

		assert.Equal(t, estr, MustSerialize(t, StripRow(t, body)))
	})

	t.Run("GET /user/find?user.like=User%25", func(t *testing.T) {
		res, err := http.Get("http://localhost:1323/user/find?name.like=User%25")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
		body := map[string]interface{}{}
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		estr := MustSerialize(t, map[string]interface{}{
			"list": []map[string]interface{}{
				{
					"Email":      "user1@foo.com",
					"ID":         2,
					"Name":       "User 1",
					"Privileges": nil,
				},
				{
					"Email":      "user2@foo.com",
					"ID":         3,
					"Name":       "User 2",
					"Privileges": nil,
				},
				{
					"Email":      "user3@foo.com",
					"ID":         4,
					"Name":       "User 3",
					"Privileges": nil,
				},
			},
			"total": 3,
		})

		assert.Equal(t, estr, MustSerialize(t, StripRow(t, body)))
	})
}

func initPrivileges(storeCtx store.StoreCtx) func(t *testing.T) {
	rs := func(s string) *string { return &s }
	return func(t *testing.T) {
		newPrivileges := []store.Privilege{
			{Name: rs("Admin"), Description: rs("Administrator")},
			{Name: rs("User"), Description: rs("User")},
			{Name: rs("Guest"), Description: rs("Guest")},
		}
		privileges, err := storeCtx.AddPrivileges(newPrivileges)
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, privileges))
		estr, _ := json.MarshalIndent([]map[string]interface{}{
			{
				"Description": "Administrator",
				"Name":        "Admin",
			},
			{
				"Description": "User",
				"Name":        "User",
			},
			{
				"Description": "Guest",
				"Name":        "Guest",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	}
}

func initUsers(storeCtx store.StoreCtx) func(t *testing.T) {
	rs := func(s string) *string { return &s }
	return func(t *testing.T) {
		newUsers := []store.User{
			{Name: rs("Admin 1"), Email: rs("admin1@cool.com"), Password: store.HashPassword("dodol123"), Privileges: &[]store.Privilege{
				{Name: rs("Admin")},
				{Name: rs("User")},
			}},
			{Name: rs("User 1"), Email: rs("user1@foo.com"), Password: store.HashPassword("dodol123"), Privileges: &[]store.Privilege{
				{Name: rs("User")},
			}},
			{Name: rs("User 2"), Email: rs("user2@foo.com"), Password: store.HashPassword("duren123"), Privileges: &[]store.Privilege{
				{Name: rs("User")},
			}},
			{Name: rs("User 3"), Email: rs("user3@foo.com"), Password: store.HashPassword("dodol123"), Privileges: &[]store.Privilege{
				{Name: rs("User")},
			}},
		}
		users, err := storeCtx.AddUsers(newUsers)
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, users))
		estr, _ := json.MarshalIndent([]map[string]interface{}{
			{
				"Name": "Admin 1",
				"Privileges": []map[string]interface{}{
					{
						"Description": "Administrator",
						"Name":        "Admin",
					},
					{
						"Description": "User",
						"Name":        "User",
					},
				},
				"Email": "admin1@cool.com",
			},
			{
				"Name": "User 1",
				"Privileges": []map[string]interface{}{
					{
						"Description": "User",
						"Name":        "User",
					},
				},
				"Email": "user1@foo.com",
			},
			{
				"Name": "User 2",
				"Privileges": []map[string]interface{}{
					{
						"Description": "User",
						"Name":        "User",
					},
				},
				"Email": "user2@foo.com",
			},
			{
				"Name": "User 3",
				"Privileges": []map[string]interface{}{
					{
						"Description": "User",
						"Name":        "User",
					},
				},
				"Email": "user3@foo.com",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	}
}
