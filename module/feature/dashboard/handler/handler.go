package handler

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"time"
)

type DashboardHandler struct {
	service dashboard.ServiceDashboardInterface
}

func NewDashboardHandler(service dashboard.ServiceDashboardInterface) dashboard.HandlerDashboardInterface {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) GetCardDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		productCount, userCount, orderCount, incomeCount, err := h.service.GetCardDashboard()
		if err != nil {
			c.Logger().Error("handler: failed to fetch dashboard data:", err.Error())
			if errors.Is(err, errors.New("gagal menghitung total produk")) {
				productCount = 0
			} else if errors.Is(err, errors.New("gagal menghitung total pelanggan")) {
				userCount = 0
			} else if errors.Is(err, errors.New("gagal menghitung total pesanan")) {
				orderCount = 0
			} else if errors.Is(err, errors.New("gagal menghitung total pendapatan")) {
				incomeCount = 0
			} else {
				return response.SendBadRequestResponse(c, "Gagal mendapatkan data: "+err.Error())
			}
		}
		return response.SendStatusOkWithDataResponse(c, "Berhasil mendapatkan data cards dasboard", dto.FormatCardResponse(productCount, userCount, orderCount, incomeCount))
	}
}

func (h *DashboardHandler) GetLandingPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		userCount, gramPlastic, orderCount, err := h.service.GetLandingPage()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data landing page")
		}

		landingPageData := dto.FormatLandingPage(userCount, gramPlastic, orderCount)
		return response.SendSuccessResponse(c, "Berhasil mendapatkan data landing page", landingPageData)
	}
}

func (h *DashboardHandler) GetReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.service.GetProductReviewsWithMaxTotal()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data ulasan")
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan data ulasan", dto.FormatLandingPageReview(result))
	}
}

func (h *DashboardHandler) GetGramPlasticStat() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		now := time.Now()
		currentMonth := now.Month()
		indonesianMonth := dto.MonthMap[currentMonth]

		startDate := time.Date(now.Year(), currentMonth, 1, 0, 0, 0, 0, now.Location())

		result := make([]dto.GramPlasticStat, 4)

		for i := 0; i < 4; i++ {
			startOfWeek := startDate.AddDate(0, 0, i*7)
			endOfWeek := startOfWeek.AddDate(0, 0, 6)

			gramTotalCount, err := h.service.GetGramPlasticStat(startOfWeek, endOfWeek)
			if err != nil {
				c.Logger().Error(fmt.Sprintf("handler: failed to fetch GramPlasticStat for week %d: %s", i+1, err.Error()))
				return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data statistik gram plastik")
			}

			result[i] = dto.GramPlasticStat{
				Week:           fmt.Sprintf("Minggu ke %d", i+1),
				GramTotalCount: gramTotalCount,
			}
		}

		return response.SendStatusOkWithDataResponses(c, "Statistik gram plastik", fmt.Sprintf("Periode bulan %s", indonesianMonth), result)
	}
}

func (h *DashboardHandler) GetLastTransactions() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		limit := 8
		transactions, err := h.service.GetLatestTransactions(limit)
		if err != nil {
			c.Logger().Error("handler: failed to fetch latest transactions:", err.Error())
			return response.SendBadRequestResponse(c, "Gagal mendapatkan data transaksi terakhir")
		}
		return response.SendStatusOkWithDataResponse(c, "Transaksi terakhir", transactions)
	}
}
