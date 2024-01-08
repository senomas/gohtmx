package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/senomas/gohtmx/store"
	"github.com/senomas/gohtmx/store/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {
	fmt.Println("TestDummy")
	storeCtx := sqlite.InitStoreCtx()
	assert.NotNil(t, storeCtx)

	newUsers := []*store.User{
		{Name: "Admin 1", Email: "admin1@foo.com", Password: store.HashPassword("dodol123")},
		{Name: "User 1", Email: "user1@foo.com", Password: store.HashPassword("dodol123")},
		{Name: "User 2", Email: "user2@foo.com", Password: store.HashPassword("duren123")},
		{Name: "User 3", Email: "user3@foo.com", Password: store.HashPassword("dodol123")},
	}
	err := storeCtx.AddUser(newUsers)
	assert.NoError(t, err)
	for _, user := range newUsers {
		json, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("ADD User: %s\n", json)
	}

	userFilter := store.UserFilter{}
	userFilter.Name.LIKE("User%")
	users, userTotal, err := storeCtx.FindUser(userFilter, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), userTotal)
	for _, user := range users {
		json, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("User: %s\n", json)
	}
}
