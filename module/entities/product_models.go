package entities

import (
	"time"
)

type ProductModels struct {
	ID            int                   `gorm:"column:id;type:int;primaryKey" json:"id"`
	Name          string                `gorm:"column:name;type:varchar(255)" json:"name"`
	Description   string                `gorm:"column:description;type:text" json:"description"`
	GramPlastic   int                   `gorm:"column:gram_plastic;type:int" json:"gram_plastic"`
	Price         float64               `gorm:"column:price;type:decimal(10,2)" json:"price"`
	Stock         int                   `gorm:"column:stock;type:int" json:"stock"`
	Discount      int                   `gorm:"column:discount;type:int" json:"discount"`
	Exp           int                   `gorm:"column:exp;type:int" json:"product_exp"`
	CreatedAt     time.Time             `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time            `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
	ProductPhotos []ProductPhotosModels `gorm:"foreignKey:ProductID" json:"product_photos"`
	ProductReview []ReviewModels        `gorm:"foreignKey:ProductID;references:ID" json:"review"`
	Categories    []CategoryModels      `gorm:"many2many:product_categories;" json:"categories"`
}

type ProductPhotosModels struct {
	ID        int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	ProductID int        `gorm:"column:product_id;type:int" json:"product_id"`
	ImageURL  string     `gorm:"column:url;type:varchar(255)" json:"url"`
	CreatedAt time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

func (ProductModels) TableName() string {
	return "products"
}

func (ProductPhotosModels) TableName() string {
	return "product_photos"
}
