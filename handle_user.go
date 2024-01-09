package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
	"github.com/senomas/gohtmx/view"
)

func wrapError(c echo.Context, vm map[string]interface{}, statusCode int, err error) {
	vm["message"] = err.Error()
	vs, _ := json.Marshal(vm)
	c.Blob(statusCode, "application/json", vs)
}

func initUserHandle(e *echo.Echo) {
	e.GET("/user/name/:name", handleUserGetByName)
	e.GET("/user/email/:email", handleUserGetByEmail)
	e.GET("/user/find", handleUserFind)
	e.GET("/user/:id", handleUserGet)
	e.POST("/user", handleUserCreate)
}

func handleUserGet(c echo.Context) error {
	appCtx := c.(*AppContext)
	storeCtx := appCtx.StoreCtx
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

func handleUserGetByName(c echo.Context) error {
	appCtx := c.(*AppContext)
	storeCtx := appCtx.StoreCtx
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

func handleUserGetByEmail(c echo.Context) error {
	appCtx := c.(*AppContext)
	storeCtx := appCtx.StoreCtx
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

func handleUserFind(c echo.Context) error {
	appCtx := c.(*AppContext)
	storeCtx := appCtx.StoreCtx
	ev := map[string]interface{}{"PARAM": c.QueryParams()}
	userFilter := store.UserFilter{}
	userFilter.Name.Set("name", c.QueryParams())
	userFilter.Email.Set("email", c.QueryParams())

	users, total, err := storeCtx.FindUsers(userFilter, 0, 100)
	if err != nil {
		wrapError(c, ev, http.StatusBadRequest, err)
		return nil
	}
	list := store.UserList{
		Users: users,
		Total: total,
	}
	accept := c.Request().Header["Accept"]
	if accept != nil && accept[0] == "application/json" {
		vs, _ := json.Marshal(list)
		return c.Blob(http.StatusOK, "application/json", vs)
	}
	return view.UserListView(list).Render(c.Request().Context(), c.Response())
}

func handleUserCreate(c echo.Context) error {
	appCtx := c.(*AppContext)
	storeCtx := appCtx.StoreCtx
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
