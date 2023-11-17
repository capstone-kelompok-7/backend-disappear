package main

import (
	"fmt"
<<<<<<< Updated upstream
	"net/http"

=======
	hAddress "github.com/capstone-kelompok-7/backend-disappear/module/feature/address/handler"
	rAddress "github.com/capstone-kelompok-7/backend-disappear/module/feature/address/repository"
	sAddress "github.com/capstone-kelompok-7/backend-disappear/module/feature/address/service"
>>>>>>> Stashed changes
	hArticle "github.com/capstone-kelompok-7/backend-disappear/module/feature/article/handler"
	rArticle "github.com/capstone-kelompok-7/backend-disappear/module/feature/article/repository"
	sArticle "github.com/capstone-kelompok-7/backend-disappear/module/feature/article/service"
	hAuth "github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/handler"
	rAuth "github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/repository"
	sAuth "github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/service"
	hCarousel "github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/handler"
	rCarousel "github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/repository"
	sCarousel "github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/service"
	hCategory "github.com/capstone-kelompok-7/backend-disappear/module/feature/category/handler"
	rCategory "github.com/capstone-kelompok-7/backend-disappear/module/feature/category/repository"
	sCategory "github.com/capstone-kelompok-7/backend-disappear/module/feature/category/service"
	hChallenge "github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/handler"
	rChallenge "github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/repository"
	sChallenge "github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/handler"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/repository"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/service"
	hVoucher "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/handler"
	rVoucher "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/repository"
	sVoucher "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/middlewares"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	hUser "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/handler"
	rUser "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/repository"
	sUser "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/service"
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
	userService := sUser.NewUserService(userRepo, hash)
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

	categoryRepo := rCategory.NewCategoryRepository(db)
	categoryService := sCategory.NewCategoryService(categoryRepo)
	categoryHandler := hCategory.NewCategoryHandler(categoryService)

	articleRepo := rArticle.NewArticleRepository(db)
	articleService := sArticle.NewArticleService(articleRepo)
	articleHandler := hArticle.NewArticleHandler(articleService)

	challengeRepo := rChallenge.NewChallengeRepository(db)
	challengeService := sChallenge.NewChallengeService(challengeRepo)
	challengeHandler := hChallenge.NewChallengeHandler(challengeService)

	carouselRepo := rCarousel.NewCarouselRepository(db)
	carouselService := sCarousel.NewCarouselService(carouselRepo)
	carouselHandler := hCarousel.NewCarouselHandler(carouselService)

	addresslRepo := rAddress.NewAddressRepository(db)
	addresslService := sAddress.NewAddressService(addresslRepo)
	addresslHandler := hAddress.NewAddressHandler(addresslService)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middlewares.ConfigureLogging())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Disappear!")
	})
	routes.RouteUser(e, userHandler, jwtService, userService)
	routes.RouteAuth(e, authHandler, jwtService, userService)
	routes.RouteVoucher(e, voucherHandler, jwtService, userService)
	routes.RouteProduct(e, productHandler, jwtService, userService)
	routes.RouteArticle(e, articleHandler, jwtService, userService)
	routes.RouteChallenge(e, challengeHandler, jwtService, userService)
	routes.RouteCategory(e, categoryHandler, jwtService, userService)
	routes.RouteCarousel(e, carouselHandler, jwtService, userService)
	routes.RouteAddress(e, addresslHandler, jwtService, userService)
	e.Logger.Fatalf(e.Start(fmt.Sprintf(":%d", initConfig.ServerPort)).Error())
}
