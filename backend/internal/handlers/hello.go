package handlers

import (
	"net/http"

	db "backend/internal/db"

	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	db.Init()
	return c.String(http.StatusOK, "Hello, World!")
}
