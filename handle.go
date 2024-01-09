package main

import (
	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
)

type AppContext struct {
	echo.Context

	StoreCtx store.StoreCtx
}

func setupRouter(e *echo.Echo) {
	initUserHandle(e)
}
