package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
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
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.AddCartItemsRequest)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		result, err := h.service.AddCartItems(currentUser.ID, req)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan produk ke keranjang: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan produk ke keranjang", result)
	}

}

func (h *CartHandler) GetCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		cartItemsSummary, err := h.service.GetCart(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail keranjang"+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail keranjang", dto.FormatCart(cartItemsSummary))
	}
}

func (h *CartHandler) ReduceQuantity() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		req := new(dto.ReduceCartItemsRequest)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		err := h.service.ReduceCartItemQuantity(req.CartItemID, req.Quantity)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengurangi kuantittas: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil mengurangi kuantittas")
	}
}

func (h *CartHandler) DeleteCartItems() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		id := c.Param("id")
		cartItemsID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		err = h.service.DeleteCartItem(cartItemsID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus produk dari keranjang: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menghapus produk dari keranjang")
	}
}
