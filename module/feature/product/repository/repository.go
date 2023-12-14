package repository

import (
	"errors"
	"strings"
	"time"

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
	query := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos").Where("deleted_at IS NULL")

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
	query := r.db.Model(&entities.ProductModels{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *ProductRepository) FindAll(page, perPage int) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos").Where("deleted_at IS NULL").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetTotalProductCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ProductModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *ProductRepository) CreateProduct(productData *entities.ProductModels, categoryIDs []uint64) (*entities.ProductModels, error) {
	tx := r.db.Begin()

	if err := tx.Create(productData).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(categoryIDs) > 0 {
		for _, categoryID := range categoryIDs {
			if err := tx.Model(productData).Association("Categories").Append(&entities.CategoryModels{ID: categoryID}); err != nil {
				tx.Rollback()
				return nil, err
			}

			if err := tx.Model(&entities.CategoryModels{}).Where("id = ?", categoryID).Update("total_product", gorm.Expr("total_product + ?", 1)).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return productData, nil
}

func (r *ProductRepository) GetProductByID(productID uint64) (*entities.ProductModels, error) {
	var products *entities.ProductModels

	if err := r.db.Preload("Categories").Preload("ProductPhotos").
		Preload("ProductReview", func(db *gorm.DB) *gorm.DB {
			return db.Limit(2)
		}).Preload("ProductReview.User").Preload("ProductReview.Photos").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&products).Error; err != nil {
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

func (r *ProductRepository) UpdateProduct(product *entities.ProductModels) (*entities.ProductModels, error) {
	tx := r.db.Begin()

	if err := tx.Save(product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan produk: " + err.Error())
	}

	return product, tx.Commit().Error
}

func (r *ProductRepository) UpdateProductCategories(product *entities.ProductModels, categoryIDs []uint64) error {
	tx := r.db.Begin()

	if err := tx.Model(product).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus asosiasi kategori: " + err.Error())
	}

	for _, categoryID := range categoryIDs {
		if err := tx.Model(product).Association("Categories").Append(&entities.CategoryModels{ID: categoryID}); err != nil {
			tx.Rollback()
			return errors.New("gagal menambahkan asosiasi kategori: " + err.Error())
		}

		if err := tx.Model(&entities.CategoryModels{}).Where("id = ?", categoryID).Update("total_product", gorm.Expr("total_product + ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("gagal memperbarui total produk kategori: " + err.Error())
		}
	}

	return tx.Commit().Error
}

func (r *ProductRepository) DeleteProduct(id uint64) error {
	var productData entities.ProductModels
	if err := r.db.First(&productData, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if err := r.db.Model(&productData).Update("DeletedAt", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProductImage(productID, imageID uint64) error {
	tx := r.db.Begin()

	if err := tx.Model(&entities.ProductPhotosModels{}).Where("id = ? AND product_id = ?", imageID, productID).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entities.ProductModels{}).Where("id = ?", productID).Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProductRepository) ReduceStockWhenPurchasing(productID, quantity uint64) error {
	var products entities.ProductModels
	if err := r.db.Model(&products).Where("id = ?", productID).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) IncreaseStock(productID, quantity uint64) error {
	var products entities.ProductModels
	if err := r.db.Model(&products).Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetTotalProductSold() (uint64, error) {
	var totalSold uint64

	query := `
        SELECT 
            COALESCE(SUM(quantity), 0) as total_sold 
        FROM 
            order_details
    `

	if err := r.db.Raw(query).Row().Scan(&totalSold); err != nil {
		return 0, err
	}

	return totalSold, nil
}

func (r *ProductRepository) GetTopRatedProducts() ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels

	if err := r.db.Order("rating desc").Where("deleted_at IS NULL").Limit(5).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductsByCategoryAndName(page, perPage int, categoryName, name string) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage

	query := r.db.Offset(offset).Limit(perPage).
		Joins("JOIN product_categories ON product_categories.product_models_id = products.id").
		Joins("JOIN category ON category.id = product_categories.category_models_id").
		Where("category.name = ? AND products.name LIKE ? AND products.deleted_at IS NULL", categoryName, "%"+name+"%")

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductsCountByCategoryAndName(categoryName, name string) (int64, error) {
	var count int64

	if err := r.db.Model(&entities.ProductModels{}).
		Joins("JOIN product_categories ON product_categories.product_models_id = products.id").
		Joins("JOIN category ON category.id = product_categories.category_models_id").
		Where("category.name = ? AND products.name LIKE ? AND products.deleted_at IS NULL", categoryName, "%"+name+"%").
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ProductRepository) GetProductsByCategoryName(categoryName string, page, perPage int) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage

	var category *entities.CategoryModels
	if err := r.db.Where("name = ?", categoryName).First(&category).Error; err != nil {
		return nil, err
	}

	if err := r.db.Offset(offset).Limit(perPage).
		Joins("JOIN product_categories ON product_categories.product_models_id = products.id").
		Joins("JOIN category ON category.id = product_categories.category_models_id").
		Where("category.name = ?", categoryName).
		Where("products.deleted_at IS NULL").
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductCountByCategoryName(categoryName string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.ProductModels{}).
		Joins("JOIN product_categories ON product_categories.product_models_id = products.id").
		Joins("JOIN category ON category.id = product_categories.category_models_id").
		Where("category.name = ?", categoryName).
		Where("products.deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ProductRepository) GetProductBySearchAndFilter(page, perPage int, sortBy, search string) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	var totalItems int64

	query := r.db.Model(&entities.ProductModels{})

	query = query.Preload("Categories").Preload("ProductPhotos").Preload("ProductReview").Where("deleted_at IS NULL")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	switch sortBy {
	case "termurah":
		query = query.Order("price asc")
	case "termahal":
		query = query.Order("price desc")
	case "abjad":
		query = query.Order("name asc")
	case "promo":
		query = query.Where("discount > ?", 0)
	case "terbaru":
		query = query.Order("created_at desc")
	}

	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *ProductRepository) GetProductByFilter(page, perPage int, sortBy string) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	var totalItems int64

	query := r.db.Model(&entities.ProductModels{})

	query = query.Preload("Categories").Preload("ProductPhotos").Preload("ProductReview").Where("deleted_at IS NULL")

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	switch sortBy {
	case "termurah":
		query = query.Order("price asc")
	case "termahal":
		query = query.Order("price desc")
	case "abjad":
		query = query.Order("name asc")
	case "promo":
		query = query.Where("discount > ?", 0)
	case "terbaru":
		query = query.Order("created_at desc")
	}

	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *ProductRepository) GetRatedProductsInRange(page, perPage int, lowerBound, upperBound float64) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage

	db := r.db.Model(&entities.ProductModels{}).Where("deleted_at IS NULL")

	err := db.Where("rating BETWEEN ? AND ?", lowerBound, upperBound).
		Order("rating").Offset(offset).Limit(perPage).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	var totalItems int64
	if err := db.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *ProductRepository) SearchByNameAndFilterByRating(page, perPage int, name, ratingParam string, lowerBound, upperBound float64) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	offset := (page - 1) * perPage
	query := r.db.Model(&entities.ProductModels{}).
		Offset(offset).
		Limit(perPage).
		Preload("Categories").
		Preload("ProductPhotos").
		Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if ratingParam != "" {
		query = query.Where("rating BETWEEN ? AND ?", lowerBound, upperBound)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *ProductRepository) FindAllProductByUserPreference(page, perPage int, productsFromAI []string) ([]*entities.ProductModels, error) {
	var matchingProducts []*entities.ProductModels

	if len(productsFromAI) == 0 {
		return nil, nil
	}

	offset := (page - 1) * perPage

	relevantQuery := r.db.Table("products")

	var relevantConditions []string
	var relevantValues []interface{}
	for _, productDesc := range productsFromAI {
		relevantConditions = append(relevantConditions, "name LIKE ?")
		relevantValues = append(relevantValues, "%"+productDesc+"%")
	}

	relevantQuery = relevantQuery.Where(strings.Join(relevantConditions, " OR "), relevantValues...)

	irrelevantQuery := r.db.Table("products").Where("deleted_at IS NULL")

	fullQuery := r.db.Raw("(?) UNION (?)", relevantQuery, irrelevantQuery).
		Preload("ProductPhotos").
		Preload("Categories").
		Limit(perPage).Offset(offset).
		Find(&matchingProducts)

	if fullQuery.Error != nil {
		return nil, fullQuery.Error
	}

	return matchingProducts, nil
}
