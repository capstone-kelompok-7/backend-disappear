package routes

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/middlewares"

	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/labstack/echo/v4"
)

func RouteAuth(e *echo.Echo, h auth.HandlerAuthInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	e.POST("api/v1/auth/register", h.Register())
	e.POST("api/v1/auth/login", h.Login())
	e.POST("/api/v1/auth/verify", h.VerifyEmail())
	e.POST("/api/v1/auth/resend-otp", h.ResendOTP())
	e.POST("/api/v1/auth/forgot-password", h.ForgotPassword())
	e.POST("/api/v1/auth/forgot-password/verify", h.VerifyOTP())
	e.POST("/api/v1/auth/forgot-password/reset", h.ResetPassword(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteUser(e *echo.Echo, h users.HandlerUserInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	users := e.Group("api/v1/users")
	users.GET("", h.GetAllUsers())
	users.GET("/by-email", h.GetUsersByEmail(), middlewares.AuthMiddleware(jwtService, userService))
	users.POST("/change-password", h.ChangePassword(), middlewares.AuthMiddleware(jwtService, userService))
	users.GET("/:id", h.GetUsersById(), middlewares.AuthMiddleware(jwtService, userService))
	users.POST("/edit-profile", h.EditProfile(), middlewares.AuthMiddleware(jwtService, userService))
	users.DELETE("/:id", h.DeleteAccount(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteVoucher(e *echo.Echo, h voucher.HandlerVoucherInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	vouchers := e.Group("api/v1/vouchers")
	vouchers.POST("", h.CreateVoucher(), middlewares.AuthMiddleware(jwtService, userService))
	vouchers.GET("", h.GetAllVouchers())
	vouchers.PUT("/:id", h.UpdateVouchers(), middlewares.AuthMiddleware(jwtService, userService))
	vouchers.GET("/:id", h.GetVoucherById())
	vouchers.DELETE("/:id", h.DeleteVoucherById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteProduct(e *echo.Echo, h product.HandlerProductInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	products := e.Group("api/v1")
	products.GET("/products", h.GetAllProducts(), middlewares.AuthMiddleware(jwtService, userService))
	products.POST("/products", h.CreateProduct(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteArticle(e *echo.Echo, h article.HandlerArticleInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	articles := e.Group("api/v1/articles")
	articles.POST("", h.CreateArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articles.GET("", h.GetAllArticles())
	articles.PUT("/:id", h.UpdateArticleById(), middlewares.AuthMiddleware(jwtService, userService))
	articles.DELETE("/:id", h.DeleteArticleById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteChallenge(e *echo.Echo, h challenge.HandlerChallengeInterface) {
	challenges := e.Group("api/v1/challenges")
	challenges.GET("", h.GetAllChallenges())
}

func RouteCategory(e *echo.Echo, h category.HandlerCategoryInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	categories := e.Group("/api/v1/categories")
	categories.POST("", h.CreateCategory(), middlewares.AuthMiddleware(jwtService, userService))
	categories.GET("", h.GetAllCategory())
	categories.GET("/:name", h.GetCategoryByName())
	categories.PUT("/:id", h.UpdateCategoryById(), middlewares.AuthMiddleware(jwtService, userService))
	categories.DELETE("/:id", h.DeleteCategoryById(), middlewares.AuthMiddleware(jwtService, userService))

}

func RouteCarousel(e *echo.Echo, h carousel.HandlerCarouselInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	carousels := e.Group("/api/v1/carousel")
	carousels.GET("", h.GetAllCarousels())
	carousels.POST("", h.CreateCarousel(), middlewares.AuthMiddleware(jwtService, userService))
	carousels.PUT("/:id", h.UpdateCarousel(), middlewares.AuthMiddleware(jwtService, userService))
	carousels.DELETE("/:id", h.DeleteCarousel(), middlewares.AuthMiddleware(jwtService, userService))
}
