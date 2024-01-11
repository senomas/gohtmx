package view

import (
	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
)

type ViewContext struct {
	echo.Context

	StoreCtx     store.StoreCtx
	userAdminCtx *UserAdminContext
}
