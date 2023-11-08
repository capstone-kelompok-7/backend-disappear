package main

import (
	"fmt"
	"net/http"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/middlewares"
	hAuth "github.com/capstone-kelompok-7/backend-disappear/module/auth/handler"
	rAuth "github.com/capstone-kelompok-7/backend-disappear/module/auth/repository"
	sAuth "github.com/capstone-kelompok-7/backend-disappear/module/auth/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/product/handler"
	"github.com/capstone-kelompok-7/backend-disappear/module/product/repository"
	"github.com/capstone-kelompok-7/backend-disappear/module/product/service"
	hUser "github.com/capstone-kelompok-7/backend-disappear/module/users/handler"
	rUser "github.com/capstone-kelompok-7/backend-disappear/module/users/repository"
	sUser "github.com/capstone-kelompok-7/backend-disappear/module/users/service"
	hVoucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/handler"
	rVoucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/repository"
	sVoucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/service"
	hArticle "github.com/capstone-kelompok-7/backend-disappear/module/article/handler"
	rArticle "github.com/capstone-kelompok-7/backend-disappear/module/article/repository"
	sArticle "github.com/capstone-kelompok-7/backend-disappear/module/article/service"
	"github.com/capstone-kelompok-7/backend-disappear/routes"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	var initConfig = config.InitConfig()

	db := database.InitDatabase(*initConfig)
	database.Migrate(db)
	jwtService := utils.NewJWT(initConfig.Secret)
	hash := utils.NewHash()

	userRepo := rUser.NewUserRepository(db)
	userService := sUser.NewUserService(userRepo)
	userHandler := hUser.NewUserHandler(userService)

	authRepo := rAuth.NewAuthRepository(db)
	authService := sAuth.NewAuthService(authRepo, jwtService, userService, hash)
	authHandler := hAuth.NewAuthHandler(authService, userService)

	voucherRepo := rVoucher.NewVoucherRepository(db)
	voucherService := sVoucher.NewVoucherService(voucherRepo)
	voucherHandler := hVoucher.NewVoucherHandler(voucherService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	articleRepo := rArticle.NewArticleRepository(db)
	articleService := sArticle.NewArticleRepository(articleRepo)
	articleHandler := hArticle.NewArticleHandler(articleService)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middlewares.ConfigureLogging())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Disappear!")
	})
	routes.RouteUser(e, userHandler)
	routes.RouteAuth(e, authHandler)
	routes.RouteVoucher(e, voucherHandler)
	routes.RouteProduct(e, productHandler)
	routes.RouteArticle(e, articleHandler)
	e.Logger.Fatalf(e.Start(fmt.Sprintf(":%d", initConfig.ServerPort)).Error())
}
