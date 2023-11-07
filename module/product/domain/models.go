package domain

import "time"

type Product struct {
	ID            int             `gorm:"column:id;type:int;primaryKey" json:"id"`
	Name          string          `gorm:"column:name;type:varchar(255)" json:"name"`
	Description   string          `gorm:"column:description;type:text" json:"description"`
	GramPlastic   int             `gorm:"column:gram_plastic;type:int" json:"gram_plastic"`
	Price         float64         `gorm:"column:price;type:decimal(10,2)" json:"price"`
	Stock         int             `gorm:"column:stock;type:int" json:"stock"`
	Discount      int             `gorm:"column:discount;type:int" json:"discount"`
	Exp           int             `gorm:"column:exp;type:int" json:"product_exp"`
	CreatedAt     time.Time       `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time      `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
	ProductPhotos []ProductPhotos `gorm:"foreignKey:ProductID" json:"product_photos"`
	ProductReview []Review        `gorm:"foreignKey:ProductID;references:ID" json:"review"`
	Categories    []Category      `gorm:"many2many:product_categories;" json:"categories"`
}

type Category struct {
	ID           int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(255)" json:"name"`
	TotalProduct int        `gorm:"column:total_product;type:int" json:"total_product"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
	Products     []Product  `gorm:"many2many:product_categories;" json:"products"`
}

type ProductPhotos struct {
	ID        int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	ProductID int        `gorm:"column:product_id;type:int" json:"product_id"`
	ImageURL  string     `gorm:"column:url;type:varchar(255)" json:"url"`
	CreatedAt time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

type Review struct {
	ID          int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	ProductID   int        `gorm:"column:product_id;type:int" json:"product_id"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Date        time.Time  `gorm:"column:date;type:date" json:"date"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}
