package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
)

func wrapError(c echo.Context, vm map[string]interface{}, statusCode int, err error) {
	vm["message"] = err.Error()
	vs, _ := json.Marshal(vm)
	c.Blob(statusCode, "application/json", vs)
}

func handleUserGet(storeCtx store.StoreCtx) func(c echo.Context) error {
	return func(c echo.Context) error {
		ids := c.Param("id")
		ev := map[string]interface{}{"PARAM_ID": ids}
		id, err := strconv.ParseInt(ids, 10, 64)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}

		user, err := storeCtx.GetUser(id)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		vs, _ := json.Marshal(user)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
}

func handleUserGetByName(storeCtx store.StoreCtx) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		ev := map[string]interface{}{"PARAM_NAME": name}
		user, err := storeCtx.GetUserByName(name)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		vs, _ := json.Marshal(user)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
}

func handleUserGetByEmail(storeCtx store.StoreCtx) func(c echo.Context) error {
	return func(c echo.Context) error {
		email := c.Param("email")
		ev := map[string]interface{}{"PARAM_EMAIL": email}
		user, err := storeCtx.GetUserByEmail(email)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		vs, _ := json.Marshal(user)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
}

func handleUserFind(storeCtx store.StoreCtx) func(c echo.Context) error {
	return func(c echo.Context) error {
		ev := map[string]interface{}{"PARAM": c.QueryParams()}
		userFilter := store.UserFilter{}
		userFilter.Name.Set("name", c.QueryParams())
		userFilter.Email.Set("email", c.QueryParams())

		users, total, err := storeCtx.FindUsers(userFilter, 0, 100)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		v := map[string]interface{}{
			"list":  users,
			"total": total,
		}
		vs, _ := json.Marshal(v)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
}

func handleUserCreate(storeCtx store.StoreCtx) func(c echo.Context) error {
	return func(c echo.Context) error {
		users := []store.User{}
		ev := map[string]interface{}{}
		err := c.Bind(&users)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		ev["BODY"] = users
		users, err = storeCtx.AddUsers(users)
		if err != nil {
			wrapError(c, ev, http.StatusBadRequest, err)
			return nil
		}
		vs, _ := json.Marshal(users)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
}
