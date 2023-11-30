package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	service dashboard.ServiceDashboardInterface
}

func NewDashboardHandler(service dashboard.ServiceDashboardInterface) dashboard.HandlerDashboardInterface {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetCardDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		productCount, userCount, orderCount, incomeCount, err := h.service.GetCardDashboard()
		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendBadRequestResponse(c, "Gagal mengambil data : "+err.Error())
		}

		return response.SendStatusOkWithDataResponse(c, "Data cards dasboard", dto.FormatCardResponse(productCount, userCount, orderCount, incomeCount))
	}
}
