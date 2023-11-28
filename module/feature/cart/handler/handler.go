package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CartHandler struct {
	service cart.ServiceCartInterface
}

func NewCartHandler(service cart.ServiceCartInterface) cart.HandlerCartInterface {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) AddCartItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.AddCartItemsRequest)
		if err := c.Bind(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		result, err := h.service.AddCartItems(currentUser.ID, req)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan ke keranjang: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil menambahkan produk ke keranjang", result)
	}

}

func (h *CartHandler) GetCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		cartItemsSummary, err := h.service.GetCart(currentUser.ID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan detail keranjang"+err.Error())
		}
		return response.SendSuccessResponse(c, "Daftar detail keranjang", dto.FormatCart(cartItemsSummary))
	}
}

func (h *CartHandler) ReduceQuantity() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		req := new(dto.ReduceCartItemsRequest)
		if err := c.Bind(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		err := h.service.ReduceCartItemQuantity(req.CartItemID, req.Quantity)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengurangi quantity: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil mengurangi quantity")
	}
}

func (h *CartHandler) DeleteCartItems() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		id := c.Param("id")
		cartItemsID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteCartItem(cartItemsID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus produk dikeranjang: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menghapus produk dikeranjang")
	}
}
