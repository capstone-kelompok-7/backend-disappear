package entities

type CartModels struct {
	ID         uint64            `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID     uint64            `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	GrandTotal uint64            `gorm:"column:grand_total;type:BIGINT UNSIGNED" json:"grand_total"`
	User       *UserModels       `json:"user" gorm:"foreignKey:UserID"`
	CartItems  []*CartItemModels `json:"cart_items,omitempty" gorm:"foreignKey:CartID"`
}

type CartItemModels struct {
	ID         uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	CartID     uint64         `gorm:"column:cart_id;type:BIGINT UNSIGNED" json:"cart_id"`
	ProductID  uint64         `gorm:"column:product_id; type:BIGINT UNSIGNED" json:"product_id"`
	Quantity   uint64         `gorm:"column:quantity;type:BIGINT UNSIGNED" json:"quantity"`
	Price      uint64         `gorm:"column:price;type:BIGINT UNSIGNED" json:"price"`
	TotalPrice uint64         `gorm:"column:total_price;type:BIGINT UNSIGNED" json:"total_price"`
	Product    *ProductModels `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (CartModels) TableName() string {
	return "carts"
}

func (CartItemModels) TableName() string {
	return "cart_items"
}
