package main

import (
	"github.com/labstack/echo/v4"
	"github.com/senomas/gohtmx/store"
)

func setupRouter(storeCtx store.StoreCtx, e *echo.Echo) {
	e.GET("/user/name/:name", handleUserGetByName(storeCtx))
	e.GET("/user/email/:email", handleUserGetByEmail(storeCtx))
	e.GET("/user/find", handleUserFind(storeCtx))
	e.GET("/user/:id", handleUserGet(storeCtx))
	e.POST("/user", handleUserCreate(storeCtx))
}
