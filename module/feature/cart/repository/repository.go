package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) cart.RepositoryCartInterface {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetCartByUserID(userID uint64) (*entities.CartModels, error) {
	carts := &entities.CartModels{}
	if err := r.db.Where("user_id = ?", userID).First(carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) GetCartByID(cartID uint64) (*entities.CartModels, error) {
	carts := &entities.CartModels{}
	if err := r.db.First(carts, cartID).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) GetCartItemByProductID(cartID, productID uint64) (*entities.CartItemModels, error) {
	var cartItem entities.CartItemModels
	if err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *CartRepository) GetCartItemsByCartID(cartID uint64) ([]*entities.CartItemModels, error) {
	var cartItems []*entities.CartItemModels
	if err := r.db.Where("cart_id = ?", cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *CartRepository) GetCartItemByID(cartItemID uint64) (*entities.CartItemModels, error) {
	var cartItem *entities.CartItemModels
	if err := r.db.Where("id = ?", cartItemID).First(&cartItem).Error; err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (r *CartRepository) GetCart(userID uint64) (*entities.CartModels, error) {
	carts := &entities.CartModels{}
	if err := r.db.
		Preload("CartItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, cart_id, product_id, quantity, total_price").
				Preload("Product", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, name, gram_plastic, price, discount").
						Preload("ProductPhotos")
				})
		}).
		Where("user_id = ?", userID).
		Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) CreateCart(newCart *entities.CartModels) (*entities.CartModels, error) {
	err := r.db.Create(newCart).Error
	if err != nil {
		return nil, err
	}
	return newCart, nil
}

func (r *CartRepository) CreateCartItem(cartItem *entities.CartItemModels) (*entities.CartItemModels, error) {
	err := r.db.Create(cartItem).Error
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (r *CartRepository) UpdateCartItem(cartItem *entities.CartItemModels) error {
	if err := r.db.Save(&cartItem).Error; err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) UpdateGrandTotal(cartID, grandTotal uint64) error {
	var carts *entities.CartModels
	result := r.db.Model(&carts).Where("id = ?", cartID).Update("grand_total", grandTotal)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CartRepository) DeleteCartItem(cartItemID uint64) error {
	result := r.db.Where("id = ?", cartItemID).Delete(&entities.CartItemModels{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CartRepository) IsProductInCart(userID, productID uint64) bool {
	var count int64
	r.db.Model(&entities.CartItemModels{}).
		Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("carts.user_id = ? AND cart_items.product_id = ?", userID, productID).
		Count(&count)
	return count > 0
}

func (r *CartRepository) RemoveProductFromCart(userID, productID uint64) error {
	var carts entities.CartModels
	if err := r.db.Where("user_id = ?", userID).Preload("CartItems").First(&carts).Error; err != nil {
		return err
	}

	var cartItem entities.CartItemModels
	for _, item := range carts.CartItems {
		if item.ProductID == productID {
			cartItem = *item
			break
		}
	}

	if cartItem.ID == 0 {
		return nil
	}

	err := r.db.Delete(&cartItem).Error
	if err != nil {
		return err
	}
	return nil
}
