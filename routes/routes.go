package routes

import (
	"github.com/capstone-kelompok-7/backend-disappear/middlewares"

	"github.com/capstone-kelompok-7/backend-disappear/module/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/labstack/echo/v4"
)

func RouteAuth(e *echo.Echo, h auth.HandlerAuthInterface) {
	e.POST("api/v1/auth/register", h.Register())
	e.POST("api/v1/auth/login", h.Login())
	e.POST("/api/v1/auth/verify", h.VerifyEmail())
}

func RouteUser(e *echo.Echo, h users.HandlerUserInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	users := e.Group("api/v1/users")
	users.GET("/list", h.GetAllUsers())
	users.GET("/by-email", h.GetUsersByEmail())
	users.PATCH("/change-password", h.ChangePassword())
	users.GET("/by-email", h.GetUsersByEmail(), middlewares.AuthMiddleware(jwtService, userService))
	users.GET("/:id", h.GetUsersById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteVoucher(e *echo.Echo, h voucher.HandlerVoucherInterface) {
	voucher := e.Group("api/v1/vouchers")
	voucher.POST("", h.CreateVoucher())
	voucher.GET("", h.GetAllVouchers())
	voucher.PUT("/:voucher_id", h.EditVoucherById())
	voucher.GET("/:voucher_id", h.GetVoucherById())
	voucher.DELETE("/:voucher_id", h.DeleteVoucherById())
}

func RouteProduct(e *echo.Echo, h product.HandlerProductInterface) {
	products := e.Group("api/v1")
	products.GET("/products", h.GetAllProducts())
}

func RouteArticle(e *echo.Echo, h article.HandlerArticleInterface) {
	articles := e.Group("api/v1/articles")
	articles.GET("", h.GetAllArticles())
}

func RouteChallenge(e *echo.Echo, h challenge.HandlerChallengeInterface) {
	challenges := e.Group("api/v1/challenges")
	challenges.GET("", h.GetAllChallenges())
}

func RouteCategory(e *echo.Echo, h category.HandlerCategoryInterface) {
	categories := e.Group("/api/v1/categories")
	categories.POST("", h.CreateCategory())
	categories.GET("", h.GetAllCategory())
	categories.GET("/:name", h.GetCategoryByName())
	categories.PUT("/:id", h.UpdateCategoryById())
	categories.DELETE("/:id", h.DeleteCategoryById())

}
