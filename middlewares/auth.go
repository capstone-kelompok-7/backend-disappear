package middlewares

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtService utils.JWTInterface, userService users.ServiceUserInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Token Bearer hilang atau tidak valid")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Token tidak valid "+err.Error())
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Token tidak valid atau telah kadaluarsa "+err.Error())
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: ID Pengguna tidak valid "+err.Error())
			}

			userID := uint64(userIDFloat)

			user, err := userService.GetUsersById(userID)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Pengguna tidak ditemukan "+err.Error())
			}

			c.Set("CurrentUser", user)

			return next(c)
		}
	}
}
