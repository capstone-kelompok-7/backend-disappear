package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.RepositoryProductInterface {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) FindByName(page, perPage int, name string) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetTotalProductCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.ProductModels{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *ProductRepository) FindAll(page, perPage int) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos").Preload("ProductReview").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetTotalProductCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ProductModels{}).Count(&count).Error
	return count, err
}

func (r *ProductRepository) CreateProduct(productData *entities.ProductModels, categoryIDs []uint64) error {
	tx := r.db.Begin()

	if err := tx.Create(productData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(categoryIDs) > 0 {
		for _, categoryID := range categoryIDs {
			if err := tx.Model(productData).Association("Categories").Append(&entities.CategoryModels{ID: categoryID}); err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Model(&entities.CategoryModels{}).Where("id = ?", categoryID).Update("total_product", gorm.Expr("total_product + ?", 1)).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *ProductRepository) GetProductByID(productID uint64) (*entities.ProductModels, error) {
	var products *entities.ProductModels

	if err := r.db.Preload("Categories").Preload("ProductPhotos").Where("id = ? AND deleted_at IS NULL", productID).First(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) CreateImageProduct(productImage *entities.ProductPhotosModels) (*entities.ProductPhotosModels, error) {
	if err := r.db.Create(&productImage).Error; err != nil {
		return productImage, err
	}
	return productImage, nil
}

func (r *ProductRepository) UpdateTotalReview(productID uint64) error {
	var products *entities.ProductModels
	err := r.db.Model(&products).Where("id = ?", productID).UpdateColumn("total_review", gorm.Expr("total_review + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) UpdateProductRating(productID uint64, newRating float64) error {
	if err := r.db.Model(&entities.ProductModels{}).Where("id = ?", productID).Update("rating", newRating).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProductReviews(page, perPage int) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage
	err := r.db.Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (r *ProductRepository) UpdateProduct(product *entities.ProductModels) error {
	tx := r.db.Begin()

	if err := tx.Save(product).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProductRepository) UpdateProductCategories(product *entities.ProductModels, categoryIDs []uint64) error {
	tx := r.db.Begin()

	if err := tx.Model(product).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	for _, categoryID := range categoryIDs {
		if err := tx.Model(product).Association("Categories").Append(&entities.CategoryModels{ID: categoryID}); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entities.CategoryModels{}).Where("id = ?", categoryID).Update("total_product", gorm.Expr("total_product + ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
