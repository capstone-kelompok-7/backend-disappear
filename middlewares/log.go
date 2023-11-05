package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureLogging() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
	})
}
