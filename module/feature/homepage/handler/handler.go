package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type HomepageHandler struct {
	service homepage.ServiceHomepageInterface
}

func NewHomepageHandler(service homepage.ServiceHomepageInterface) homepage.HandlerHomepageInterface {
	return &HomepageHandler{
		service: service,
	}
}

func (h *HomepageHandler) GetHomepageContent() echo.HandlerFunc {
	return func(c echo.Context) error {
		limit := 5
		bestSellingProducts, err := h.service.GetBestSellingProducts(limit)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan produk terlaris")
		}
		categories, err := h.service.GetCategory()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan kategori")
		}
		carousel, err := h.service.GetCarousel()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan carousel")
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan konten halaman utama", dto.FormatContentResponse(carousel, categories, bestSellingProducts))
	}
}

func (h *HomepageHandler) GetBlogPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		challenges, err := h.service.GetChallenge()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan challenge")
		}
		articles, err := h.service.GetArticle()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan artikel")
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan konten halaman utama", dto.FormatBlogPostResponse(challenges, articles))
	}
}
