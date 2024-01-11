package view

import (
	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/stores"
)

type ViewContext struct {
	echo.Context

	store        stores.Store
	userAdminCtx *UserAdminContext
}
