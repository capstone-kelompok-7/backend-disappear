package routes

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/labstack/echo/v4"
)

func RouteAuth(e *echo.Echo, h auth.HandlerAuthInterface) {
	e.POST("api/v1/auth/register", h.Register())
	e.POST("api/v1/auth/login", h.Login())
}

func RouteUser(e *echo.Echo, h users.HandlerUserInterface) {
	users := e.Group("api/v1/users")
	users.GET("/list", h.GetAllUsers())
	users.GET("/by-email", h.GetUsersByEmail())
}

func RouteVoucher(e *echo.Echo, h voucher.HandlerVoucherInterface) {
	voucher := e.Group("api/v1/vouchers")
	voucher.POST("", h.CreateVoucher())
	voucher.GET("", h.GetAllVouchers())
}
