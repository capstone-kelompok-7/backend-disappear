package routes

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/middlewares"
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
	authGroup.POST("/register-social", h.RegisterSocial())
	authGroup.POST("/login-social", h.LoginSocial())
}

func RouteUser(e *echo.Echo, h users.HandlerUserInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	usersGroup := e.Group("api/v1/users")
	usersGroup.GET("", h.GetAllUsers(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/by-email", h.GetUsersByEmail(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.POST("/change-password", h.ChangePassword(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/:id", h.GetUsersById(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.POST("/edit-profile", h.EditProfile(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.DELETE("/:id", h.DeleteAccount(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/leaderboard", h.GetLeaderboard(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/get-activities/:id", h.GetUserTransactionActivity(), middlewares.AuthMiddleware(jwtService, userService))
	usersGroup.GET("/get-profile", h.GetUserProfile(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteVoucher(e *echo.Echo, h voucher.HandlerVoucherInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	voucherGroup := e.Group("api/v1/vouchers")
	voucherGroup.POST("", h.CreateVoucher(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("", h.GetAllVouchers(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.PUT("/:id", h.UpdateVouchers(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("/:id", h.GetVoucherById(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.DELETE("/:id", h.DeleteVoucherById(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.POST("/claims", h.ClaimVoucher(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("/users", h.GetVoucherUser(), middlewares.AuthMiddleware(jwtService, userService))
	voucherGroup.GET("/to-claims", h.GetAllVouchersToClaims(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteProduct(e *echo.Echo, h product.HandlerProductInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	productsGroup := e.Group("api/v1/products")
	productsGroup.GET("", h.GetAllProducts())
	productsGroup.POST("", h.CreateProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/:id", h.GetProductById())
	productsGroup.POST("/images", h.CreateProductImage(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/reviews", h.GetAllProductsReview(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.PUT("/:id", h.UpdateProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.DELETE("/:id", h.DeleteProduct(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.DELETE("/:idProduct/image/:idImage", h.DeleteProductImageById(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/preferences", h.GetAllProductsPreferences(), middlewares.AuthMiddleware(jwtService, userService))
	productsGroup.GET("/other-products", h.GetTopRatedProducts(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteArticle(e *echo.Echo, h article.HandlerArticleInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	articlesGroup := e.Group("api/v1/articles")
	articlesGroup.POST("", h.CreateArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("", h.GetAllArticles(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("/:id", h.GetArticleById(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.PUT("/:id", h.UpdateArticleById(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.DELETE("/:id", h.DeleteArticleById(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.POST("/bookmark", h.BookmarkArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.DELETE("/bookmark/:id", h.DeleteBookmarkedArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("/bookmark", h.GetUsersBookmark(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("/preferences", h.GetAllArticleUser(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("/other-article", h.GetOtherArticle(), middlewares.AuthMiddleware(jwtService, userService))
	articlesGroup.GET("/latest-article", h.GetLatestArticle(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteChallenge(e *echo.Echo, h challenge.HandlerChallengeInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	challengesGroup := e.Group("api/v1/challenges")
	challengesGroup.GET("", h.GetAllChallenges(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.POST("", h.CreateChallenge(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.PUT("/:id", h.UpdateChallenge(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.DELETE("/:id", h.DeleteChallengeById(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.GET("/:id", h.GetChallengeById(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.POST("/submit", h.CreateSubmitChallengeForm(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.PUT("/participants/status/:id", h.UpdateSubmitChallengeForm(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.GET("/participants", h.GetAllSubmitChallengeForm(), middlewares.AuthMiddleware(jwtService, userService))
	challengesGroup.GET("/participants/:id", h.GetSubmitChallengeFormById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteCategory(e *echo.Echo, h category.HandlerCategoryInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	categoriesGroup := e.Group("/api/v1/categories")
	categoriesGroup.GET("", h.GetAllCategory(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.GET("/:id", h.GetCategoryById(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.POST("", h.CreateCategory(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.PUT("/:id", h.UpdateCategoryById(), middlewares.AuthMiddleware(jwtService, userService))
	categoriesGroup.DELETE("/:id", h.DeleteCategoryById(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteCarousel(e *echo.Echo, h carousel.HandlerCarouselInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	carouselsGroup := e.Group("/api/v1/carousel")
	carouselsGroup.GET("", h.GetAllCarousels(), middlewares.AuthMiddleware(jwtService, userService))
	carouselsGroup.GET("/:id", h.GetCarouselById(), middlewares.AuthMiddleware(jwtService, userService))
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

func RouteCart(e *echo.Echo, h cart.HandlerCartInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	cartGroup := e.Group("/api/v1/carts")
	cartGroup.POST("", h.AddCartItem(), middlewares.AuthMiddleware(jwtService, userService))
	cartGroup.GET("", h.GetCart(), middlewares.AuthMiddleware(jwtService, userService))
	cartGroup.PUT("/reduce/quantity", h.ReduceQuantity(), middlewares.AuthMiddleware(jwtService, userService))
	cartGroup.DELETE("/cart-items/:id", h.DeleteCartItems(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteOrder(e *echo.Echo, h order.HandlerOrderInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	orderGroup := e.Group("/api/v1/order")
	orderGroup.GET("", h.GetAllOrders(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.GET("/payment", h.GetAllPayment(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.GET("/:id", h.GetOrderById(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.POST("", h.CreateOrder(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.POST("/confirm/:id", h.ConfirmPayment(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.POST("/carts", h.CreateOrderFromCart(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.POST("/cancel/:id", h.CancelPayment(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.POST("/callback", h.Callback())
	orderGroup.PUT("/update-order", h.UpdateOrderStatus(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.GET("/by-users", h.GetAllOrderByUserID(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.PUT("/accept-order/:id", h.AcceptOrder(), middlewares.AuthMiddleware(jwtService, userService))
	orderGroup.GET("/track", h.Tracking(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteAssistant(e *echo.Echo, h assistant.HandlerAssistantInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	assistantGroup := e.Group("/api/v1/assistant")
	assistantGroup.POST("/question", h.CreateQuestion(), middlewares.AuthMiddleware(jwtService, userService))
	assistantGroup.POST("/answer", h.CreateAnswer(), middlewares.AuthMiddleware(jwtService, userService))
	assistantGroup.GET("", h.GetChatByIdUser(), middlewares.AuthMiddleware(jwtService, userService))
	assistantGroup.POST("/generate-article", h.GenerateArticle(), middlewares.AuthMiddleware(jwtService, userService))
	assistantGroup.GET("/product", h.GetProductByIdUser(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteDashboard(e *echo.Echo, h dashboard.HandlerDashboardInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	dashboardGroup := e.Group("/api/v1/dashboards")
	dashboardGroup.GET("/card", h.GetCardDashboard(), middlewares.AuthMiddleware(jwtService, userService))
	dashboardGroup.GET("/landing-page", h.GetLandingPage())
	dashboardGroup.GET("/reviews", h.GetReview())
	dashboardGroup.GET("/chart", h.GetGramPlasticStat(), middlewares.AuthMiddleware(jwtService, userService))
	dashboardGroup.GET("/transactions", h.GetLastTransactions(), middlewares.AuthMiddleware(jwtService, userService))

}

func RouteHomepage(e *echo.Echo, h homepage.HandlerHomepageInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	homeGroup := e.Group("/api/v1/homepage")
	homeGroup.GET("/blog-posts", h.GetBlogPost(), middlewares.AuthMiddleware(jwtService, userService))
	homeGroup.GET("/content", h.GetHomepageContent(), middlewares.AuthMiddleware(jwtService, userService))
}

func RouteFcm(e *echo.Echo, h fcm.HandlerFcmInterface, jwtService utils.JWTInterface, userService users.ServiceUserInterface) {
	fcmGroup := e.Group("/api/v1/fcm")
	fcmGroup.GET("/:id", h.GetFcmById(), middlewares.AuthMiddleware(jwtService, userService))
	fcmGroup.GET("/users", h.GetFcmByIdUser(), middlewares.AuthMiddleware(jwtService, userService))
	fcmGroup.POST("", h.CreateFcm(), middlewares.AuthMiddleware(jwtService, userService))
	fcmGroup.PUT("/:id", h.DeleteFcmById(), middlewares.AuthMiddleware(jwtService, userService))
}
