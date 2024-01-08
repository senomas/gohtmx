package store

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/senomas/gohtmx/store"
	_ "github.com/senomas/gohtmx/store/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestUserPrivilegeCrud(t *testing.T) {
	db_type := os.Getenv("DB_TYPE")
	if db_type == "" {
		db_type = "sqlite"
	}
	storeCtx := store.Get(db_type)
	assert.NotNil(t, storeCtx)
	defer storeCtx.Close()

	rs := func(s string) *string { return &s }

	t.Run("initialize privilege", func(t *testing.T) {
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
	})

	var user1ID int64

	t.Run("initialize user", func(t *testing.T) {
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
			{Name: rs("Demo"), Email: rs("demo@foo.com"), Password: store.HashPassword("dodol123"), Privileges: &[]store.Privilege{}},
		}
		users, err := storeCtx.AddUsers(newUsers)
		for _, user := range users {
			if *user.Name == "User 1" {
				user1ID = user.ID
			}
		}
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
			{
				"Name":       "Demo",
				"Privileges": []map[string]interface{}{},
				"Email":      "demo@foo.com",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("Get user user1 and verify password", func(t *testing.T) {
		user, err := storeCtx.GetUser(user1ID)
		assert.NoError(t, err)
		assert.Equal(t, "User 1", *user.Name)
		assert.Truef(t, store.VerifyPassword("dodol123", *user.Password), "invalid password")
		str := MustSerialize(t, StripRow(t, user))
		estr, _ := json.MarshalIndent(map[string]interface{}{
			"Name": "User 1",
			"Privileges": []map[string]interface{}{
				{
					"Description": "User",
					"Name":        "User",
				},
			},
			"Email": "user1@foo.com",
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("Get demo", func(t *testing.T) {
		user, err := storeCtx.GetUserByName("Demo")
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, user))
		estr, _ := json.MarshalIndent(map[string]interface{}{
			"Name":       "Demo",
			"Privileges": []map[string]interface{}{},
			"Email":      "demo@foo.com",
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("Update demo set privileges", func(t *testing.T) {
		user, err := storeCtx.GetUserByName("Demo")
		assert.NoError(t, err)
		updateUser := store.User{
			ID: user.ID,
			Privileges: &[]store.Privilege{
				{Name: rs("Admin")},
				{Name: rs("User")},
				{Name: rs("Guest")},
			},
		}
		err = storeCtx.UpdateUser(updateUser)
		assert.NoError(t, err)
		user, err = storeCtx.GetUser(user.ID)
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, user))
		estr, _ := json.MarshalIndent(map[string]interface{}{
			"Name": "Demo",
			"Privileges": []map[string]interface{}{
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
			},
			"Email": "demo@foo.com",
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("Update demo set less privileges", func(t *testing.T) {
		user, err := storeCtx.GetUserByName("Demo")
		assert.NoError(t, err)
		updateUser := store.User{
			ID: user.ID,
			Privileges: &[]store.Privilege{
				{Name: rs("User")},
			},
		}
		err = storeCtx.UpdateUser(updateUser)
		assert.NoError(t, err)
		user, err = storeCtx.GetUser(user.ID)
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, user))
		estr, _ := json.MarshalIndent(map[string]interface{}{
			"Name": "Demo",
			"Privileges": []map[string]interface{}{
				{
					"Description": "User",
					"Name":        "User",
				},
			},
			"Email": "demo@foo.com",
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("Update demo name", func(t *testing.T) {
		user, err := storeCtx.GetUserByName("Demo")
		assert.NoError(t, err)
		updateUser := store.User{
			ID:    user.ID,
			Name:  rs("User 4"),
			Email: rs("demo-user4@foo.com"),
		}
		err = storeCtx.UpdateUser(updateUser)
		assert.NoError(t, err)
		user, err = storeCtx.GetUser(user.ID)
		assert.NoError(t, err)
		str := MustSerialize(t, StripRow(t, user))
		estr, _ := json.MarshalIndent(map[string]interface{}{
			"Name": "User 4",
			"Privileges": []map[string]interface{}{
				{
					"Description": "User",
					"Name":        "User",
				},
			},
			"Email": "demo-user4@foo.com",
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("find name like User% include demo", func(t *testing.T) {
		userFilter := store.UserFilter{}
		userFilter.Name.Like("User%")
		users, userTotal, err := storeCtx.FindUsers(userFilter, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), userTotal)
		str := MustSerialize(t, StripRow(t, users))
		estr, _ := json.MarshalIndent([]map[string]interface{}{
			{
				"Name":       "User 1",
				"Privileges": nil,
				"Email":      "user1@foo.com",
			},
			{
				"Name":       "User 2",
				"Privileges": nil,
				"Email":      "user2@foo.com",
			},
			{
				"Name":       "User 3",
				"Privileges": nil,
				"Email":      "user3@foo.com",
			},
			{
				"Name":       "User 4",
				"Privileges": nil,
				"Email":      "demo-user4@foo.com",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("delete user email demo-user4", func(t *testing.T) {
		user, err := storeCtx.GetUserByEmail("demo-user4@foo.com")
		assert.NoError(t, err)
		err = storeCtx.DeleteUsers([]int64{user.ID})
		assert.NoError(t, err)

		t.Run("check user_privileges", func(t *testing.T) {
			up, err := storeCtx.GetUserPrivileges(user.ID)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(up))
		})
	})

	t.Run("find name like User%", func(t *testing.T) {
		userFilter := store.UserFilter{}
		userFilter.Name.Like("User%")
		users, userTotal, err := storeCtx.FindUsers(userFilter, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), userTotal)
		str := MustSerialize(t, StripRow(t, users))
		estr, _ := json.MarshalIndent([]map[string]interface{}{
			{
				"Name":       "User 1",
				"Privileges": nil,
				"Email":      "user1@foo.com",
			},
			{
				"Name":       "User 2",
				"Privileges": nil,
				"Email":      "user2@foo.com",
			},
			{
				"Name":       "User 3",
				"Privileges": nil,
				"Email":      "user3@foo.com",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("find email like %@cool.com", func(t *testing.T) {
		userFilter := store.UserFilter{}
		userFilter.Email.Like("%@cool.com")
		users, userTotal, err := storeCtx.FindUsers(userFilter, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), userTotal)
		str := MustSerialize(t, StripRow(t, users))
		estr, _ := json.MarshalIndent([]map[string]interface{}{
			{
				"Name":       "Admin 1",
				"Privileges": nil,
				"Email":      "admin1@cool.com",
			},
		}, "", "  ")
		assert.Equal(t, string(estr), string(str))
	})

	t.Run("delete used privilege", func(t *testing.T) {
		privilege, err := storeCtx.GetPrivilegeByName("User")
		assert.NoError(t, err)
		err = storeCtx.DeletePrivileges([]int64{privilege.ID})
		assert.ErrorContains(t, err, "record in use")
	})
}

func StripRow(t *testing.T, row interface{}) interface{} {
	rv := reflect.ValueOf(row)
	if rv.Kind() == reflect.Slice {
		res := []interface{}{}
		for i := 0; i < rv.Len(); i++ {
			vv := rv.Index(i).Interface()
			res = append(res, StripRow(t, vv))
		}
		return res
	}
	if rv.Kind() == reflect.Struct {
		tv := rv.Type()
		res := map[string]interface{}{}
		for i := 0; i < tv.NumField(); i++ {
			k := tv.Field(i).Name
			switch k {
			case "ID":
				// skip
			case "Password":
				// skip
			default:
				v := rv.Field(i).Interface()
				res[k] = StripRow(t, v)
			}
		}
		return res
	}
	if rv.Kind() == reflect.Map {
		res := map[string]interface{}{}
		for _, k := range rv.MapKeys() {
			switch k.String() {
			case "ID":
				// skip
			case "Password":
				// skip
			default:
				vv := rv.MapIndex(k).Interface()
				res[k.String()] = StripRow(t, vv)
			}
		}
		return res
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		ve := rv.Elem()
		if ve.IsValid() {
			return StripRow(t, ve.Interface())
		}
	}
	return row
}

func MustSerialize(t *testing.T, v interface{}) string {
	str, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return string(str)
}
