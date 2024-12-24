package contract

import "github.com/labstack/echo/v4"

type Handler interface {
	Handler(next echo.HandlerFunc) echo.HandlerFunc
}
