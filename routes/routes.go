package routes

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/middlewares"

	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/labstack/echo/v4"
)

func RouteAuth(e *echo.Echo, h auth.HandlerAuthInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	authGroup := e.Group("api/v1/auth")
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/verify", h.VerifyEmail())
	authGroup.POST("/resend-otp", h.ResendOTP())
	authGroup.POST("/forgot-password", h.ForgotPassword())
	authGroup.POST("/forgot-password/verify", h.VerifyOTP())
	authGroup.POST("/forgot-password/reset", h.ResetPassword(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteUser(e *echo.Echo, h users.HandlerUserInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	usersGroup := e.Group("api/v1/users")
	usersGroup.GET("", h.GetAllUsers())
	usersGroup.GET("/by-email", h.GetUsersByEmail(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.POST("/change-password", h.ChangePassword(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/:id", h.GetUsersById(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.POST("/edit-profile", h.EditProfile(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.DELETE("/:id", h.DeleteAccount(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteVoucher(e *echo.Echo, h voucher.HandlerVoucherInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	voucherGroup := e.Group("api/v1/vouchers")
	voucherGroup.POST("", h.CreateVoucher(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("", h.GetAllVouchers())
	voucherGroup.PUT("/:id", h.UpdateVouchers(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("/:id", h.GetVoucherById())
	voucherGroup.DELETE("/:id", h.DeleteVoucherById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteProduct(e *echo.Echo, h product.HandlerProductInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	productsGroup := e.Group("api/v1/products")
	productsGroup.GET("", h.GetAllProducts(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.POST("", h.CreateProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/:id", h.GetProductById())
	productsGroup.POST("/images", h.CreateProductImage(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/reviews", h.GetAllProductsReview(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.PUT("/:id", h.UpdateProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.DELETE("/:id", h.DeleteProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.DELETE("/:idProduct/image/:idImage", h.DeleteProductImageById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteArticle(e *echo.Echo, h article.HandlerArticleInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	articlesGroup := e.Group("api/v1/articles")
	articlesGroup.POST("", h.CreateArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("", h.GetAllArticles())
	articlesGroup.PUT("/:id", h.UpdateArticleById(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.DELETE("/:id", h.DeleteArticleById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteChallenge(e *echo.Echo, h challenge.HandlerChallengeInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	challengesGroup := e.Group("api/v1/challenges")
	challengesGroup.GET("", h.GetAllChallenges())
	challengesGroup.POST("", h.CreateChallenge(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.PUT("/:id", h.UpdateChallenge(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.DELETE("/:id", h.DeleteChallengeById(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.GET("/:id", h.GetChallengeById())
}

func RouteCategory(e *echo.Echo, h category.HandlerCategoryInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	categoriesGroup := e.Group("/api/v1/categories")
	categoriesGroup.POST("", h.CreateCategory(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.GET("", h.GetAllCategory())
	categoriesGroup.GET("/:name", h.GetCategoryByName())
	categoriesGroup.PUT("/:id", h.UpdateCategoryById(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.DELETE("/:id", h.DeleteCategoryById(), middlewares.AuthMiddleware(jwtService, userService))

}

func RouteCarousel(e *echo.Echo, h carousel.HandlerCarouselInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	carouselsGroup := e.Group("/api/v1/carousel")
	carouselsGroup.GET("", h.GetAllCarousels())
	carouselsGroup.POST("", h.CreateCarousel(), middlewares.AuthMiddleware(jwtService, userService))
	carouselsGroup.PUT("/:id", h.UpdateCarousel(), middlewares.AuthMiddleware(jwtService, userService))
	carouselsGroup.DELETE("/:id", h.DeleteCarousel(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteAddress(e *echo.Echo, h address.HandlerAddressInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	addressesGroup := e.Group("/api/v1/address")
	addressesGroup.GET("", h.GetAllAddress(), middlewares.AuthMiddleware(jwtService, userService))
	addressesGroup.POST("", h.CreateAddress(), middlewares.AuthMiddleware(jwtService, userService))
	addressesGroup.PUT("/:id", h.UpdateAddress(), middlewares.AuthMiddleware(jwtService, userService))
	addressesGroup.DELETE("/:id", h.DeleteAddress(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteReview(e *echo.Echo, h review.HandlerReviewInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	reviewsGroup := e.Group("/api/v1/reviews")
	reviewsGroup.POST("", h.CreateReview(), middlewares.AuthMiddleware(jwtService, userService))
	reviewsGroup.POST("/photos", h.CreateReviewImages(), middlewares.AuthMiddleware(jwtService, userService))
	reviewsGroup.GET("/:id", h.GetReviewById(), middlewares.AuthMiddleware(jwtService, userService))
	reviewsGroup.GET("/detail/:id", h.GetDetailReviewProduct(), middlewares.AuthMiddleware(jwtService, userService))
}
