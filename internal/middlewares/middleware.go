package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CoockieHash(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		coockie := new(http.Cookie)
		coockie.Name = "user"
		coockie.Value = "31337"
		c.SetCookie(coockie)
		return next(c)
	}
}
