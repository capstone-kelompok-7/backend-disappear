package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

// OrderResponse OrderResponse Respon Get Order By ID
type OrderResponse struct {
	ID                    string                `json:"id"`
	IdOrder               string                `json:"id_order"`
	AddressID             uint64                `json:"address_id"`
	UserID                uint64                `json:"user_id"`
	VoucherID             uint64                `json:"voucher_id"`
	Note                  string                `json:"note"`
	GrandTotalGramPlastic uint64                `json:"grand_total_gram_plastic"`
	GrandTotalExp         uint64                `json:"grand_total_exp"`
	GrandTotalQuantity    uint64                `json:"grand_total_quantity"`
	GrandTotalPrice       uint64                `json:"grand_total_price"`
	ShipmentFee           uint64                `json:"shipment_fee"`
	AdminFees             uint64                `json:"admin_fees"`
	GrandTotalDiscount    uint64                `json:"grand_total_discount"`
	TotalAmountPaid       uint64                `json:"total_amount_paid"`
	OrderStatus           string                `json:"order_status"`
	PaymentStatus         string                `json:"payment_status"`
	PaymentMethod         string                `json:"payment_method"`
	ExtraInfo             string                `json:"extra_info"`
	StatusOrderDate       time.Time             `json:"status_order_date"`
	CreatedAt             time.Time             `json:"created_at"`
	Address               AddressResponse       `json:"address"`
	User                  UserResponse          `json:"user"`
	Voucher               VoucherResponse       `json:"voucher"`
	OrderDetails          []OrderDetailResponse `json:"order_details"`
}

type OrderDetailResponse struct {
	ID               uint64          `json:"id"`
	OrderID          string          `json:"order_id"`
	ProductID        uint64          `json:"product_id"`
	Quantity         uint64          `json:"quantity"`
	TotalGramPlastic uint64          `json:"total_gram_plastic"`
	TotalExp         uint64          `json:"total_exp"`
	TotalPrice       uint64          `json:"total_price"`
	TotalDiscount    uint64          `json:"total_discount"`
	Product          ProductResponse `json:"product,omitempty"`
}

type ProductPhotoResponse struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	URL       string `json:"url"`
}

type ProductResponse struct {
	ID            uint64                 `json:"id"`
	Name          string                 `json:"name"`
	Price         uint64                 `json:"price"`
	Discount      uint64                 `json:"discount"`
	GramPlastic   uint64                 `json:"gram_plastic"`
	ProductExp    uint64                 `json:"product_exp"`
	ProductPhotos []ProductPhotoResponse `json:"product_photos"`
}

type AddressResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	AcceptedName string `json:"accepted_name" `
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	IsPrimary    bool   `json:"is_primary"`
}

type UserResponse struct {
	ID           uint64 `json:"id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	PhotoProfile string `json:"photo_profile"`
}

type VoucherResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Discount    uint64 `json:"discount"`
	MinPurchase uint64 `json:"min_purchase"`
}

func FormatOrderDetail(order *entities.OrderModels) OrderResponse {
	orderResponse := OrderResponse{
		ID:                    order.ID,
		IdOrder:               order.IdOrder,
		AddressID:             order.AddressID,
		UserID:                order.UserID,
		VoucherID:             0,
		Note:                  order.Note,
		GrandTotalGramPlastic: order.GrandTotalGramPlastic,
		GrandTotalExp:         order.GrandTotalExp,
		GrandTotalQuantity:    order.GrandTotalQuantity,
		GrandTotalPrice:       order.GrandTotalPrice,
		ShipmentFee:           order.ShipmentFee,
		AdminFees:             order.AdminFees,
		GrandTotalDiscount:    order.GrandTotalDiscount,
		TotalAmountPaid:       order.TotalAmountPaid,
		OrderStatus:           order.OrderStatus,
		PaymentStatus:         order.PaymentStatus,
		PaymentMethod:         order.PaymentMethod,
		StatusOrderDate:       order.StatusOrderDate,
		ExtraInfo:             order.ExtraInfo,
		CreatedAt:             order.CreatedAt,
		Address: AddressResponse{
			ID:           order.Address.ID,
			UserID:       order.Address.UserID,
			AcceptedName: order.Address.AcceptedName,
			Phone:        order.Address.Phone,
			Address:      order.Address.Address,
			IsPrimary:    order.Address.IsPrimary,
		},
		User: UserResponse{
			ID:           order.User.ID,
			Email:        order.User.Email,
			Phone:        order.User.Phone,
			Name:         order.User.Name,
			PhotoProfile: order.User.PhotoProfile,
		},
		Voucher: VoucherResponse{
			ID:          order.Voucher.ID,
			Name:        order.Voucher.Name,
			Code:        order.Voucher.Code,
			Category:    order.Voucher.Category,
			Description: order.Voucher.Description,
			Discount:    order.Voucher.Discount,
			MinPurchase: order.Voucher.MinPurchase,
		},
	}

	var orderDetails []OrderDetailResponse
	for _, detail := range order.OrderDetails {
		var productPhotos []ProductPhotoResponse
		for _, photo := range detail.Product.ProductPhotos {
			productPhotos = append(productPhotos, ProductPhotoResponse{
				ID:        photo.ID,
				ProductID: photo.ProductID,
				URL:       photo.ImageURL,
			})
		}

		orderDetail := OrderDetailResponse{
			ID:               detail.ID,
			OrderID:          detail.OrderID,
			ProductID:        detail.ProductID,
			Quantity:         detail.Quantity,
			TotalGramPlastic: detail.TotalGramPlastic,
			TotalExp:         detail.TotalExp,
			TotalPrice:       detail.TotalPrice,
			TotalDiscount:    detail.TotalDiscount,
			Product: ProductResponse{
				ID:            detail.Product.ID,
				Name:          detail.Product.Name,
				Price:         detail.Product.Price,
				Discount:      detail.Product.Discount,
				GramPlastic:   detail.Product.GramPlastic,
				ProductExp:    detail.Product.Exp,
				ProductPhotos: productPhotos,
			},
		}
		if len(detail.Product.ProductPhotos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:        detail.Product.ProductPhotos[0].ID,
				ProductID: detail.Product.ProductPhotos[0].ProductID,
				URL:       detail.Product.ProductPhotos[0].ImageURL,
			}
			orderDetail.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		orderDetails = append(orderDetails, orderDetail)
	}
	if order.VoucherID != nil {
		orderResponse.VoucherID = *order.VoucherID
	}

	orderResponse.OrderDetails = orderDetails

	return orderResponse
}

// OrderPaginationResponse Pagination Response
type OrderPaginationResponse struct {
	ID              string                      `json:"id"`
	IdOrder         string                      `json:"id_order"`
	UserID          uint64                      `json:"user_id"`
	TotalAmountPaid uint64                      `json:"total_amount_paid"`
	OrderStatus     string                      `json:"order_status"`
	CreatedAt       time.Time                   `json:"created_at"`
	User            UserPaginationOrderResponse `json:"user"`
}

type UserPaginationOrderResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func FormatOrderPagination(order *entities.OrderModels) *OrderPaginationResponse {
	orderResponse := &OrderPaginationResponse{
		ID:              order.ID,
		IdOrder:         order.IdOrder,
		UserID:          order.UserID,
		TotalAmountPaid: order.TotalAmountPaid,
		OrderStatus:     order.OrderStatus,
		CreatedAt:       order.CreatedAt,
		User: UserPaginationOrderResponse{
			ID:   order.User.ID,
			Name: order.User.Name,
		},
	}
	return orderResponse
}

func FormatterOrder(orders []*entities.OrderModels) []*OrderPaginationResponse {
	var orderFormatters []*OrderPaginationResponse

	for _, order := range orders {
		formattedOrder := FormatOrderPagination(order)
		orderFormatters = append(orderFormatters, formattedOrder)
	}

	return orderFormatters
}

// PaymentPaginationResponse Pagination Response
type PaymentPaginationResponse struct {
	ID              string                      `json:"id"`
	IdOrder         string                      `json:"id_order"`
	UserID          uint64                      `json:"user_id"`
	TotalAmountPaid uint64                      `json:"total_amount_paid"`
	PaymentStatus   string                      `json:"payment_status"`
	CreatedAt       time.Time                   `json:"created_at"`
	User            UserPaginationOrderResponse `json:"user"`
}

type UserPaymentPaginationOrderResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func FormatOrderPaymentPagination(order *entities.OrderModels) *PaymentPaginationResponse {
	orderResponse := &PaymentPaginationResponse{
		ID:              order.ID,
		IdOrder:         order.IdOrder,
		UserID:          order.UserID,
		TotalAmountPaid: order.TotalAmountPaid,
		PaymentStatus:   order.PaymentStatus,
		CreatedAt:       order.CreatedAt,
		User: UserPaginationOrderResponse{
			ID:   order.User.ID,
			Name: order.User.Name,
		},
	}
	return orderResponse
}

func FormatterOrderPayment(orders []*entities.OrderModels) []*PaymentPaginationResponse {
	var orderFormatters []*PaymentPaginationResponse

	for _, order := range orders {
		formattedOrder := FormatOrderPaymentPagination(order)
		orderFormatters = append(orderFormatters, formattedOrder)
	}

	return orderFormatters
}

// OrderCreationResponse Create Response
type OrderCreationResponse struct {
	ID                    string                        `json:"id"`
	IdOrder               string                        `json:"id_order"`
	AddressID             uint64                        `json:"address_id"`
	UserID                uint64                        `json:"user_id"`
	VoucherID             uint64                        `json:"voucher_id"`
	Note                  string                        `json:"note"`
	GrandTotalGramPlastic uint64                        `json:"grand_total_gram_plastic"`
	GrandTotalExp         uint64                        `json:"grand_total_exp"`
	GrandTotalQuantity    uint64                        `json:"grand_total_quantity"`
	GrandTotalPrice       uint64                        `json:"grand_total_price"`
	ShipmentFee           uint64                        `json:"shipment_fee"`
	AdminFees             uint64                        `json:"admin_fees"`
	GrandTotalDiscount    uint64                        `json:"grand_total_discount"`
	TotalAmountPaid       uint64                        `json:"total_amount_paid"`
	OrderStatus           string                        `json:"order_status"`
	PaymentStatus         string                        `json:"payment_status"`
	PaymentMethod         string                        `json:"payment_method"`
	CreatedAt             time.Time                     `json:"created_at"`
	OrderDetails          []OrderDetailCreationResponse `json:"order_details"`
}

type OrderDetailCreationResponse struct {
	ID               uint64 `json:"id"`
	OrderID          string `json:"order_id"`
	ProductID        uint64 `json:"product_id"`
	Quantity         uint64 `json:"quantity"`
	TotalGramPlastic uint64 `json:"total_gram_plastic"`
	TotalExp         uint64 `json:"total_exp"`
	TotalPrice       uint64 `json:"total_price"`
	TotalDiscount    uint64 `json:"total_discount"`
}

func CreateOrderFormatter(order *entities.OrderModels) OrderCreationResponse {
	orderResponse := OrderCreationResponse{
		ID:                    order.ID,
		IdOrder:               order.IdOrder,
		AddressID:             order.AddressID,
		UserID:                order.UserID,
		VoucherID:             0,
		Note:                  order.Note,
		GrandTotalGramPlastic: order.GrandTotalGramPlastic,
		GrandTotalExp:         order.GrandTotalExp,
		GrandTotalQuantity:    order.GrandTotalQuantity,
		GrandTotalPrice:       order.GrandTotalPrice,
		ShipmentFee:           order.ShipmentFee,
		AdminFees:             order.AdminFees,
		GrandTotalDiscount:    order.GrandTotalDiscount,
		TotalAmountPaid:       order.TotalAmountPaid,
		OrderStatus:           order.OrderStatus,
		PaymentStatus:         order.PaymentStatus,
		PaymentMethod:         order.PaymentMethod,
		CreatedAt:             order.CreatedAt,
	}

	var orderDetails []OrderDetailCreationResponse
	for _, detail := range order.OrderDetails {
		orderDetail := OrderDetailCreationResponse{
			ID:               detail.ID,
			OrderID:          detail.OrderID,
			ProductID:        detail.ProductID,
			Quantity:         detail.Quantity,
			TotalGramPlastic: detail.TotalGramPlastic,
			TotalExp:         detail.TotalExp,
			TotalPrice:       detail.TotalPrice,
			TotalDiscount:    detail.TotalDiscount,
		}
		orderDetails = append(orderDetails, orderDetail)
	}
	if order.VoucherID != nil {
		orderResponse.VoucherID = *order.VoucherID
	}
	orderResponse.OrderDetails = orderDetails
	return orderResponse
}

// GetAllOrderUserResponse Respon to Get Order By UserID
type GetAllOrderUserResponse struct {
	ID              string                `json:"id"`
	IdOrder         string                `json:"id_order"`
	UserID          uint64                `json:"user_id"`
	Note            string                `json:"note"`
	TotalAmountPaid uint64                `json:"total_amount_paid"`
	OrderStatus     string                `json:"order_status"`
	PaymentStatus   string                `json:"payment_status"`
	PaymentMethod   string                `json:"payment_method"`
	ExtraInfo       string                `json:"extra_info"`
	StatusOrderDate time.Time             `json:"status_order_date"`
	CreatedAt       time.Time             `json:"created_at"`
	OrderDetails    []OrderDetailResponse `json:"order_details"`
}

func FormatGetAllOrderUser(order *entities.OrderModels) *GetAllOrderUserResponse {
	orderResponse := &GetAllOrderUserResponse{
		ID:              order.ID,
		IdOrder:         order.IdOrder,
		UserID:          order.UserID,
		Note:            order.Note,
		TotalAmountPaid: order.TotalAmountPaid,
		OrderStatus:     order.OrderStatus,
		PaymentStatus:   order.PaymentStatus,
		PaymentMethod:   order.PaymentMethod,
		StatusOrderDate: order.StatusOrderDate,
		ExtraInfo:       order.ExtraInfo,
		CreatedAt:       order.CreatedAt,
	}

	var orderDetails []OrderDetailResponse
	for _, detail := range order.OrderDetails {
		var productPhotos []ProductPhotoResponse
		for _, photo := range detail.Product.ProductPhotos {
			productPhotos = append(productPhotos, ProductPhotoResponse{
				ID:        photo.ID,
				ProductID: photo.ProductID,
				URL:       photo.ImageURL,
			})
		}

		orderDetail := OrderDetailResponse{
			ID:               detail.ID,
			OrderID:          detail.OrderID,
			ProductID:        detail.ProductID,
			Quantity:         detail.Quantity,
			TotalGramPlastic: detail.TotalGramPlastic,
			TotalExp:         detail.TotalExp,
			TotalPrice:       detail.TotalPrice,
			TotalDiscount:    detail.TotalDiscount,
			Product: ProductResponse{
				ID:            detail.Product.ID,
				Name:          detail.Product.Name,
				Price:         detail.Product.Price,
				Discount:      detail.Product.Discount,
				GramPlastic:   detail.Product.GramPlastic,
				ProductExp:    detail.Product.Exp,
				ProductPhotos: productPhotos,
			},
		}
		if len(detail.Product.ProductPhotos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:        detail.Product.ProductPhotos[0].ID,
				ProductID: detail.Product.ProductPhotos[0].ProductID,
				URL:       detail.Product.ProductPhotos[0].ImageURL,
			}
			orderDetail.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	orderResponse.OrderDetails = orderDetails

	return orderResponse
}

func FormatterGetAllOrderUser(orders []*entities.OrderModels) []*GetAllOrderUserResponse {
	var orderFormatters []*GetAllOrderUserResponse

	for _, order := range orders {
		formattedOrder := FormatGetAllOrderUser(order)
		orderFormatters = append(orderFormatters, formattedOrder)
	}

	return orderFormatters
}
