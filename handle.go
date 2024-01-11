package main

import (
	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/view"
)

func setupRouter(e *echo.Echo) {
	e.GET("/user", view.UserAdminHandler)
}
