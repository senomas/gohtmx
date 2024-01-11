package view

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
)

type UserAdminContext struct {
	Error  error
	Edit   *store.User
	View   *store.User
	List   []store.User
	Filter store.UserFilter
	Offset int64
	Limit  int
	Total  int64
}

func UserAdminHandler(c echo.Context) error {
	viewCtx := c.(*ViewContext)
	storeCtx := viewCtx.StoreCtx

	pc := viewCtx.userAdminCtx
	if pc == nil {
		pc = &UserAdminContext{}
		pc.Offset = 0
		pc.Total = 10
		viewCtx.userAdminCtx = pc
	}
	v := c.QueryParam("_o")
	if v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			pc.Offset = i
		}
	}
	v = c.QueryParam("_l")
	if v != "" {
		i, err := strconv.ParseInt(v, 10, 8)
		if err != nil && i >= 5 && i <= 100 {
			pc.Limit = int(i)
		}
	}
	pc.Filter.Name.Set("name", c.QueryParams())
	pc.Filter.Email.Set("email", c.QueryParams())
	users, total, err := storeCtx.FindUsers(pc.Filter, pc.Offset, pc.Limit)
	if err != nil {
		pc.Error = err
	} else {
		pc.List = users
		pc.Total = total
	}

	return UserAdmin(pc).Render(c.Request().Context(), c.Response())
}
