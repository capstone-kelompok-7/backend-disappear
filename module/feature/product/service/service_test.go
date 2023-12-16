package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	repoAi "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/mocks"
	serviceAi "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/mocks"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestProductService_CalculatePaginationValues(t *testing.T) {
	service := &ProductService{}

	t.Run("Page less than or equal to zero should default to 1", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(0, 100, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page exceeds total pages should set to total pages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(15, 100, 8)

		assert.Equal(t, 13, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page within limits should return correct values", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(2, 100, 8)

		assert.Equal(t, 2, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Total items not perfectly divisible by perPage should round totalPages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(1, 95, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 12, totalPages)
	})
}

func TestProductService_GetNextPage(t *testing.T) {
	service := &ProductService{}

	t.Run("Next Page Within Total Pages", func(t *testing.T) {
		currentPage := 3
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, currentPage+1, nextPage)
	})

	t.Run("Next Page Equal to Total Pages", func(t *testing.T) {
		currentPage := 5
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, totalPages, nextPage)
	})
}

func TestProductService_GetPrevPage(t *testing.T) {
	service := &ProductService{}

	t.Run("Previous Page Within Bounds", func(t *testing.T) {
		currentPage := 3

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage-1, prevPage)
	})

	t.Run("Previous Page at Lower Bound", func(t *testing.T) {
		currentPage := 1

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage, prevPage)
	})
}

func TestProductService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatBot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatBot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	products := []*entities.ProductModels{
		{
			ID:          1,
			Name:        "Product test 1",
			Description: "Description 1",
			GramPlastic: 1,
			Price:       100,
			Stock:       10,
			Discount:    2,
			Exp:         3,
			Rating:      4,
			TotalReview: 1,
		},
		{
			ID:          2,
			Name:        "Product test 2",
			Description: "Description 2",
			GramPlastic: 1,
			Price:       100,
			Stock:       10,
			Discount:    2,
			Exp:         3,
			Rating:      4,
			TotalReview: 1,
		},
	}

	t.Run("Succes case - Product Found", func(t *testing.T) {
		expectedTotalItems := int64(8)
		repo.On("FindAll", 1, 8).Return(products, nil).Once()
		repo.On("GetTotalProductCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAll(1, 8)
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetTotalProductCount Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalProductCount Error")
		repo.On("FindAll", 1, 8).Return(nil, nil).Once()
		repo.On("GetTotalProductCount").Return(int64(0), expectedErr).Once()

		product, totalItems, err := service.GetAll(1, 8)

		assert.Nil(t, product)
		assert.Error(t, err)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		repo.On("FindAll", 1, 8).Return(nil, expectedErr).Once()

		product, totalItems, err := service.GetAll(1, 8)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)

	})
}

func TestProductService_GetProductsByName(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatBot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatBot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	products := []*entities.ProductModels{
		{
			ID:          1,
			Name:        "Product test 1",
			Description: "Description 1",
			GramPlastic: 1,
			Price:       100,
			Stock:       10,
			Discount:    2,
			Exp:         3,
			Rating:      4,
			TotalReview: 1,
		},
		{
			ID:          2,
			Name:        "Product test 2",
			Description: "Description 2",
			GramPlastic: 1,
			Price:       100,
			Stock:       10,
			Discount:    2,
			Exp:         3,
			Rating:      4,
			TotalReview: 1,
		},
	}
	name := "Product Test"
	t.Run("Succes Case - Product Found by Name", func(t *testing.T) {
		expectedTotalItems := int64(8)
		repo.On("FindByName", 1, 8, name).Return(products, nil).Once()
		repo.On("GetTotalProductCountByName", name).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetProductsByName(1, 8, name)
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Product by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to find product by name")
		repo.On("FindByName", 1, 8, name).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetProductsByName(1, 8, name)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Product Count by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get total product count by name")
		repo.On("FindByName", 1, 8, name).Return(products, nil).Once()
		repo.On("GetTotalProductCountByName", name).Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetProductsByName(1, 8, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatBot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatBot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	request := &dto.CreateProductRequest{
		Name:        "Product Test",
		Description: "Description Test",
		GramPlastic: 1,
		Price:       100,
		Stock:       10,
		Discount:    2,
		Exp:         3,
		ImageURL:    "https://example.com/image.jpg",
	}

	t.Run("Success case - Product Created", func(t *testing.T) {
		createdProduct := &entities.ProductModels{
			ID:          1,
			Name:        request.Name,
			Description: request.Description,
			GramPlastic: request.GramPlastic,
			Price:       request.Price,
			Stock:       request.Stock,
			Discount:    request.Discount,
			Exp:         request.Exp,
		}

		repo.On("CreateProduct", mock.Anything, mock.Anything).Return(createdProduct, nil).Once()

		result, err := service.CreateProduct(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdProduct.ID, result.ID)
		assert.Equal(t, createdProduct.Name, result.Name)
		assert.Equal(t, createdProduct.Description, result.Description)
		assert.Equal(t, createdProduct.GramPlastic, result.GramPlastic)
		assert.Equal(t, createdProduct.Price, result.Price)
		assert.Equal(t, createdProduct.Stock, result.Stock)
		assert.Equal(t, createdProduct.Discount, result.Discount)
		assert.Equal(t, createdProduct.Exp, result.Exp)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - CreateProduct Error", func(t *testing.T) {
		expectedErr := errors.New("CreateProduct Error")
		repo.On("CreateProduct", mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

		result, err := service.CreateProduct(request)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Success case - Product Created with ImageURL", func(t *testing.T) {
		createdProduct := &entities.ProductModels{
			ID:   1,
			Name: request.Name,
			ProductPhotos: []entities.ProductPhotosModels{
				{
					ImageURL:  request.ImageURL,
					CreatedAt: time.Now(),
				},
			},
		}

		repo.On("CreateProduct", mock.Anything, mock.Anything).Return(createdProduct, nil).Once()

		result, err := service.CreateProduct(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdProduct.ID, result.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Success case - Product Created without ImageURL", func(t *testing.T) {
		createdProduct := &entities.ProductModels{
			ID:   2,
			Name: request.Name,
		}

		repo.On("CreateProduct", mock.Anything, mock.Anything).Return(createdProduct, nil).Once()

		request.ImageURL = ""

		result, err := service.CreateProduct(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdProduct.ID, result.ID)
		assert.Len(t, result.ProductPhotos, 0)

		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductByID(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatBot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatBot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	product := &entities.ProductModels{
		ID:          1,
		Name:        "Product test 1",
		Description: "Description 1",
		GramPlastic: 1,
		Price:       100,
		Stock:       10,
		Discount:    2,
		Exp:         3,
		Rating:      4,
		TotalReview: 1,
	}

	expectedProduct := &entities.ProductModels{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		GramPlastic: product.GramPlastic,
		Price:       product.Price,
		Stock:       product.Stock,
		Discount:    product.Discount,
		Exp:         product.Exp,
		Rating:      product.Rating,
		TotalReview: product.TotalReview,
	}

	t.Run("Succes Case - Product found", func(t *testing.T) {
		productId := uint64(1)
		repo.On("GetProductByID", productId).Return(expectedProduct, nil).Once()

		result, err := service.GetProductByID(productId)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct.ID, result.ID)
		assert.Equal(t, expectedProduct.Name, result.Name)
		assert.Equal(t, expectedProduct.Description, result.Description)
		assert.Equal(t, expectedProduct.GramPlastic, result.GramPlastic)
		assert.Equal(t, expectedProduct.Price, result.Price)
		assert.Equal(t, expectedProduct.Stock, result.Stock)
		assert.Equal(t, expectedProduct.Discount, result.Discount)
		assert.Equal(t, expectedProduct.Exp, result.Exp)
		assert.Equal(t, expectedProduct.Rating, result.Rating)
		assert.Equal(t, expectedProduct.TotalReview, result.TotalReview)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Product Not Found", func(t *testing.T) {
		productId := uint64(2)
		expectedErr := errors.New("produk tidak ditemukan")
		repo.On("GetProductByID", productId).Return(nil, expectedErr).Once()

		result, err := service.GetProductByID(productId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})
}

func TestProductService_CreateImageProduct(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	request := dto.CreateProductImage{
		ProductID: 1,
		Image:     "https://example.com/image.jpg",
	}

	t.Run("Success case - Image Created for Product", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:          1,
			Name:        "Product Test",
			Description: "Description Test",
		}

		repo.On("GetProductByID", request.ProductID).Return(product, nil).Once()

		createdImage := &entities.ProductPhotosModels{
			ProductID: product.ID,
			ImageURL:  request.Image,
		}

		repo.On("CreateImageProduct", createdImage).Return(createdImage, nil).Once()

		result, err := service.CreateImageProduct(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdImage.ProductID, result.ProductID)
		assert.Equal(t, createdImage.ImageURL, result.ImageURL)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductByID Error", func(t *testing.T) {
		expectedErr := errors.New("GetProductByID Error")
		repo.On("GetProductByID", request.ProductID).Return(nil, expectedErr).Once()

		result, err := service.CreateImageProduct(request)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - CreateImageProduct Error", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: 1,
		}

		repo.On("GetProductByID", request.ProductID).Return(product, nil).Once()

		expectedErr := errors.New("CreateImageProduct Error")
		repo.On("CreateImageProduct", mock.Anything).Return(nil, expectedErr).Once()

		result, err := service.CreateImageProduct(request)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})
}

func TestProductService_UpdateTotalReview(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)

	t.Run("Success case - Total Review Updated", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:          productID,
			Name:        "Product Test",
			Description: "Description Test",
			TotalReview: 10,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("UpdateTotalReview", productID).Return(nil).Once()

		err := service.UpdateTotalReview(productID)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductByID Error", func(t *testing.T) {
		expectedErr := errors.New("GetProductByID Error")
		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()

		err := service.UpdateTotalReview(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - UpdateTotalReview Error", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: productID,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		expectedErr := errors.New("UpdateTotalReview Error")
		repo.On("UpdateTotalReview", productID).Return(expectedErr).Once()

		err := service.UpdateTotalReview(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal memperbarui total review")

		repo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProductRating(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)
	newRating := 4.5

	t.Run("Success case - Product Rating Updated", func(t *testing.T) {
		repo.On("UpdateProductRating", productID, newRating).Return(nil).Once()

		err := service.UpdateProductRating(productID, newRating)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - UpdateProductRating Error", func(t *testing.T) {
		expectedErr := errors.New("UpdateProductRating Error")
		repo.On("UpdateProductRating", productID, newRating).Return(expectedErr).Once()

		err := service.UpdateProductRating(productID, newRating)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal memperbarui rating produk")

		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductReviews(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	page := 1
	perPage := 10

	t.Run("Success case - Get Product Reviews", func(t *testing.T) {
		products := []*entities.ProductModels{
			{
				ID:   1,
				Name: "Product Test 1",
			},
			{
				ID:   2,
				Name: "Product Test 2",
			},
		}

		expectedTotalItems := int64(20)

		repo.On("GetProductReviews", page, perPage).Return(products, nil).Once()
		repo.On("GetTotalProductCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetProductReviews(page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(products), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductReviews Error", func(t *testing.T) {
		expectedErr := errors.New("GetProductReviews Error")

		repo.On("GetProductReviews", page, perPage).Return(nil, expectedErr).Once()

		result, _, err := service.GetProductReviews(page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetTotalProductCount Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalProductCount Error")

		repo.On("GetProductReviews", page, perPage).Return(nil, nil).Once()
		repo.On("GetTotalProductCount").Return(int64(0), expectedErr).Once()

		result, _, err := service.GetProductReviews(page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)

	updateRequest := &dto.UpdateProduct{
		Name:        "Updated Product",
		Description: "Updated Description",
		GramPlastic: 2,
		Price:       200,
		Stock:       20,
		Discount:    5,
		Exp:         6,
		ImageURL:    "https://example.com/updated_image.jpg",
		CategoryIDs: []uint64{1, 2},
	}

	t.Run("Success case - Product Updated", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:          productID,
			Name:        "Product Test",
			Description: "Description Test",
			ProductPhotos: []entities.ProductPhotosModels{
				{
					ImageURL:  "https://example.com/original_image.jpg",
					CreatedAt: time.Now(),
				},
			},
			Categories: []entities.CategoryModels{
				{
					ID:   1,
					Name: "Category 1",
				},
			},
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		updateRequest := &dto.UpdateProduct{
			Name:        "Updated Product",
			Description: "Updated Description",
			GramPlastic: 2,
			Price:       200,
			Stock:       20,
			Discount:    5,
			Exp:         6,
			ImageURL:    "https://example.com/updated_image.jpg",
			CategoryIDs: []uint64{2, 3},
		}

		repo.On("UpdateProductCategories", product, updateRequest.CategoryIDs).Return(nil).Once()

		repo.On("UpdateProduct", product).Return(product, nil).Once()

		updatedProduct, err := service.UpdateProduct(productID, updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, updateRequest.Name, updatedProduct.Name)
		assert.Equal(t, updateRequest.Description, updatedProduct.Description)
		assert.Equal(t, updateRequest.GramPlastic, updatedProduct.GramPlastic)
		assert.Equal(t, updateRequest.Price, updatedProduct.Price)
		assert.Equal(t, updateRequest.Stock, updatedProduct.Stock)
		assert.Equal(t, updateRequest.Discount, updatedProduct.Discount)
		assert.Equal(t, updateRequest.Exp, updatedProduct.Exp)
		assert.Equal(t, updateRequest.ImageURL, updatedProduct.ProductPhotos[0].ImageURL)
		assert.Equal(t, 1, len(updatedProduct.Categories))

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductByID Error", func(t *testing.T) {
		expectedErr := errors.New("GetProductByID Error")
		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()

		result, err := service.UpdateProduct(productID, updateRequest)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - UpdateProductCategories Error", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: productID,
		}

		expectedErr := errors.New("UpdateProductCategories Error")
		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("UpdateProductCategories", product, updateRequest.CategoryIDs).Return(expectedErr).Once()

		result, err := service.UpdateProduct(productID, updateRequest)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - UpdateProduct Error", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: productID,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("UpdateProductCategories", product, updateRequest.CategoryIDs).Return(nil).Once()

		expectedErr := errors.New("UpdateProduct Error")
		repo.On("UpdateProduct", product).Return(nil, expectedErr).Once()

		result, err := service.UpdateProduct(productID, updateRequest)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)

	t.Run("Success case - Product Deleted", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: productID,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		repo.On("DeleteProduct", productID).Return(nil).Once()

		err := service.DeleteProduct(productID)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Product Not Found", func(t *testing.T) {
		expectedErr := errors.New("produk tidak ditemukan")

		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()

		err := service.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - DeleteProduct Error", func(t *testing.T) {
		expectedErr := errors.New("gagal menghapus product")

		product := &entities.ProductModels{
			ID: productID,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		repo.On("DeleteProduct", productID).Return(expectedErr).Once()

		err := service.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})
}

func TestProductService_DeleteImageProduct(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)
	imageID := uint64(1)

	t.Run("Success case - Image Deleted", func(t *testing.T) {
		product := &entities.ProductModels{
			ID: productID,
			ProductPhotos: []entities.ProductPhotosModels{
				{
					ID:        imageID,
					ImageURL:  "https://example.com/image.jpg",
					CreatedAt: time.Now(),
				},
			},
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("DeleteProductImage", productID, imageID).Return(nil).Once()

		err := service.DeleteImageProduct(productID, imageID)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Product Not Found", func(t *testing.T) {
		expectedErr := errors.New("produk tidak ditemukan")

		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()

		err := service.DeleteImageProduct(productID, imageID)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Image Not Found", func(t *testing.T) {
		expectedErr := errors.New("image tidak ditemukan pada produk ini")

		product := &entities.ProductModels{
			ID: productID,
			ProductPhotos: []entities.ProductPhotosModels{
				{
					ID:        2,
					ImageURL:  "https://example.com/other_image.jpg",
					CreatedAt: time.Now(),
				},
			},
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		err := service.DeleteImageProduct(productID, imageID)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - DeleteProductImage Error", func(t *testing.T) {
		expectedErr := errors.New("gagal menghapus image pada produk ini")

		product := &entities.ProductModels{
			ID: productID,
			ProductPhotos: []entities.ProductPhotosModels{
				{
					ID:        imageID,
					ImageURL:  "https://example.com/image.jpg",
					CreatedAt: time.Now(),
				},
			},
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()

		repo.On("DeleteProductImage", productID, imageID).Return(expectedErr).Once()

		err := service.DeleteImageProduct(productID, imageID)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())

		repo.AssertExpectations(t)
	})
}

func TestProductService_ReduceStockWhenPurchasing(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)
	quantity := uint64(5)

	t.Run("Success case - Stock Reduced", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:    productID,
			Stock: 10,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("ReduceStockWhenPurchasing", productID, quantity).Return(nil).Once()

		err := service.ReduceStockWhenPurchasing(productID, quantity)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Product Not Found", func(t *testing.T) {
		expectedErr := errors.New("produk tidak ditemukan")
		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()

		err := service.ReduceStockWhenPurchasing(productID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Insufficient Stock", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:    productID,
			Stock: 3,
		}
		repo.On("GetProductByID", productID).Return(product, nil).Once()

		err := service.ReduceStockWhenPurchasing(productID, quantity)
		assert.Error(t, err)
		assert.EqualError(t, err, "stok tidak mencukupi untuk pesanan ini")
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - ReduceStockWhenPurchasing Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mengurangi stok")

		product := &entities.ProductModels{
			ID:    productID,
			Stock: 10,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("ReduceStockWhenPurchasing", productID, quantity).Return(expectedErr).Once()

		err := service.ReduceStockWhenPurchasing(productID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})
}

func TestProductService_IncreaseStock(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	productID := uint64(1)
	quantity := uint64(5)

	t.Run("Success case - Stock Increased", func(t *testing.T) {
		product := &entities.ProductModels{
			ID:    productID,
			Stock: 10,
		}

		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("IncreaseStock", productID, quantity).Return(nil).Once()

		err := service.IncreaseStock(productID, quantity)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Product Not Found", func(t *testing.T) {
		expectedErr := errors.New("produk tidak ditemukan")

		repo.On("GetProductByID", productID).Return(nil, expectedErr).Once()
		err := service.IncreaseStock(productID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - IncreaseStock Error", func(t *testing.T) {
		expectedErr := errors.New("gagal menambah stok")

		product := &entities.ProductModels{
			ID:    productID,
			Stock: 10,
		}
		repo.On("GetProductByID", productID).Return(product, nil).Once()
		repo.On("IncreaseStock", productID, quantity).Return(expectedErr).Once()

		err := service.IncreaseStock(productID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetTotalProductSold(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	t.Run("Success case - Total Product Sold", func(t *testing.T) {
		repo.On("GetTotalProductSold").Return(uint64(100), nil).Once()

		totalSold, err := service.GetTotalProductSold()

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), totalSold)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetTotalProductSold Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total produk terjual")

		repo.On("GetTotalProductSold").Return(uint64(0), expectedErr).Once()

		totalSold, err := service.GetTotalProductSold()

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Equal(t, uint64(0), totalSold)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetTopRatedProducts(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoAi := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoAi, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	t.Run("Success case - Top Rated Products", func(t *testing.T) {
		topRatedProduct := &entities.ProductModels{
			ID:          1,
			Name:        "Top Rated Product",
			Description: "Description",
		}

		repo.On("GetTopRatedProducts").Return([]*entities.ProductModels{topRatedProduct}, nil).Once()

		result, err := service.GetTopRatedProducts()

		assert.NoError(t, err)
		assert.Equal(t, []*entities.ProductModels{topRatedProduct}, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetTopRatedProducts Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan produk dengan rating tertinggi")

		repo.On("GetTopRatedProducts").Return(nil, expectedErr).Once()

		result, err := service.GetTopRatedProducts()

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductsByCategoryAndName(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	categoryName := "Electronics"
	name := "Smartphone"
	page := 1
	perPage := 10

	t.Run("Success case - Products Found", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:          1,
				Name:        "Smartphone A",
				Description: "Description A",
			},
			{
				ID:          2,
				Name:        "Smartphone B",
				Description: "Description B",
			},
		}

		repo.On("GetProductsByCategoryAndName", page, perPage, categoryName, name).Return(foundProducts, nil).Once()
		repo.On("GetProductsCountByCategoryAndName", categoryName, name).Return(int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.GetProductsByCategoryAndName(categoryName, name, page, perPage)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductsByCategoryAndName Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan produk sesuai pencarian")

		repo.On("GetProductsByCategoryAndName", page, perPage, categoryName, name).Return(nil, expectedErr).Once()

		result, _, err := service.GetProductsByCategoryAndName(categoryName, name, page, perPage)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductsCountByCategoryAndName Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan jumlah total produk sesuai pencarian")
		repo.On("GetProductsByCategoryAndName", page, perPage, categoryName, name).Return([]*entities.ProductModels{}, nil).Once()
		repo.On("GetProductsCountByCategoryAndName", categoryName, name).Return(int64(0), expectedErr).Once()

		result, _, err := service.GetProductsByCategoryAndName(categoryName, name, page, perPage)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductsByCategoryName(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	categoryName := "Electronics"
	page := 1
	perPage := 10

	t.Run("Success case - Products Found", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:          1,
				Name:        "Product A",
				Description: "Description A",
			},
			{
				ID:          2,
				Name:        "Product B",
				Description: "Description B",
			},
		}
		repo.On("GetProductsByCategoryName", categoryName, page, perPage).Return(foundProducts, nil).Once()
		repo.On("GetProductCountByCategoryName", categoryName).Return(int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.GetProductsByCategoryName(categoryName, page, perPage)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductsByCategoryName Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan produk sesuai pencarian")

		repo.On("GetProductsByCategoryName", categoryName, page, perPage).Return(nil, expectedErr).Once()

		result, _, err := service.GetProductsByCategoryName(categoryName, page, perPage)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductCountByCategoryName Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan jumlah total produk sesuai pencarian")

		repo.On("GetProductsByCategoryName", categoryName, page, perPage).Return([]*entities.ProductModels{}, nil).Once()
		repo.On("GetProductCountByCategoryName", categoryName).Return(int64(0), expectedErr).Once()

		result, _, err := service.GetProductsByCategoryName(categoryName, page, perPage)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductsBySearchAndFilter(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	page := 1
	perPage := 8
	filter := "abjad"
	search := "smartphone"

	t.Run("Success case - Products Found with Filter", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:          1,
				Name:        "Smartphone A",
				Description: "Description A",
				Categories: []entities.CategoryModels{
					{
						ID:   1,
						Name: "smartphone",
					},
				},
			},
			{
				ID:          2,
				Name:        "Smartphone B",
				Description: "Description B",
			},
		}
		repo.On("GetProductBySearchAndFilter", page, perPage, filter, search).Return(foundProducts, int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.GetProductsBySearchAndFilter(page, perPage, filter, search)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Success case - Products Found without Filter", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:          3,
				Name:        "Laptop C",
				Description: "Description C",
			},
			{
				ID:          4,
				Name:        "Tablet D",
				Description: "Description D",
			},
		}

		repo.On("GetProductBySearchAndFilter", page, perPage, "", search).Return(foundProducts, int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.GetProductsBySearchAndFilter(page, perPage, "", search)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed case - GetProductBySearchAndFilter Error", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan produk sesuai pencarian dan filter")

		repo.On("GetProductBySearchAndFilter", page, perPage, mock.Anything, search).Return(nil, int64(0), expectedErr).Once()

		result, _, err := service.GetProductsBySearchAndFilter(page, perPage, filter, search)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestProductService_GetProductsByFilter(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatBot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatBot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	page := 1
	perPage := 10
	filter := "abjad"

	t.Run("Success Case - Valid Filter", func(t *testing.T) {
		expectedProducts := []*entities.ProductModels{}
		expectedTotalItems := int64(len(expectedProducts))
		repo.On("GetProductByFilter", page, perPage, filter).Return(expectedProducts, expectedTotalItems, nil).Once()

		products, totalItems, err := service.GetProductsByFilter(page, perPage, filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")
		repo.On("GetProductByFilter", page, perPage, filter).Return(nil, int64(0), expectedErr).Once()

		products, totalItems, err := service.GetProductsByFilter(page, perPage, filter)
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestGetRatingBounds(t *testing.T) {
	t.Run("Success Case - Valid Rating", func(t *testing.T) {
		ratingParam := "sangat buruk"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, 0.1, lowerBound)
		assert.Equal(t, 1.9, upperBound)
	})

	t.Run("Success Case - Valid Rating", func(t *testing.T) {
		ratingParam := "buruk"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, float64(2), lowerBound)
		assert.Equal(t, 2.9, upperBound)
	})

	t.Run("Success Case - Valid Rating", func(t *testing.T) {
		ratingParam := "sedang"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, float64(3), lowerBound)
		assert.Equal(t, 3.9, upperBound)
	})

	t.Run("Success Case - Valid Rating", func(t *testing.T) {
		ratingParam := "baik"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, float64(4), lowerBound)
		assert.Equal(t, 4.9, upperBound)
	})

	t.Run("Success Case - Valid Rating", func(t *testing.T) {
		ratingParam := "sangat baik"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, float64(5), lowerBound)
		assert.Equal(t, 5.0, upperBound)
	})

	t.Run("Error Case - Invalid Rating", func(t *testing.T) {
		ratingParam := "unknown"
		lowerBound, upperBound, err := getRatingBounds(ratingParam)

		assert.Error(t, err)
		assert.Zero(t, lowerBound)
		assert.Zero(t, upperBound)
		assert.Equal(t, errors.New("tipe filter tidak valid"), err)
	})
}

func TestProductService_GetRatedProductsInRange(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	page := 1
	perPage := 8
	ratingParam := "baik"

	t.Run("Success case - Rated Products Found", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:     1,
				Name:   "Product A",
				Rating: 4.5,
			},
			{
				ID:     2,
				Name:   "Product B",
				Rating: 4.8,
			},
		}

		expectedLowerBound, expectedUpperBound, _ := getRatingBounds(ratingParam)
		repo.On("GetRatedProductsInRange", page, perPage, expectedLowerBound, expectedUpperBound).Return(foundProducts, int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.GetRatedProductsInRange(page, perPage, ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Error case - Invalid Rating Param", func(t *testing.T) {
		ratingParam := "invalid"

		result, totalItems, err := service.GetRatedProductsInRange(page, perPage, ratingParam)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
	})

	t.Run("Error case - Repository error", func(t *testing.T) {
		expectedLowerBound, expectedUpperBound, _ := getRatingBounds(ratingParam)

		repo.On("GetRatedProductsInRange", page, perPage, expectedLowerBound, expectedUpperBound).Return(nil, int64(0), errors.New("some repository error")).Once()

		result, totalItems, err := service.GetRatedProductsInRange(page, perPage, ratingParam)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})
}

func TestProductService_SearchByNameAndFilterByRating(t *testing.T) {
	repo := mocks.NewRepositoryProductInterface(t)
	repoChatbot := repoAi.NewRepositoryAssistantInterface(t)
	var initConfig = config.InitConfig()
	var client = openai.NewClient(initConfig.OpenAiApiKey)
	serviceAI := serviceAi.NewAssistantService(repoChatbot, client, *initConfig)
	service := NewProductService(repo, serviceAI)

	page := 1
	perPage := 8
	name := "Product"
	ratingParam := "baik"

	t.Run("Success case - Products Found", func(t *testing.T) {
		foundProducts := []*entities.ProductModels{
			{
				ID:     1,
				Name:   "Product A",
				Rating: 4.5,
			},
			{
				ID:     2,
				Name:   "Product B",
				Rating: 4.8,
			},
		}

		expectedLowerBound, expectedUpperBound, _ := getRatingBounds(ratingParam)
		repo.On("SearchByNameAndFilterByRating", page, perPage, name, ratingParam, expectedLowerBound, expectedUpperBound).Return(foundProducts, int64(len(foundProducts)), nil).Once()

		result, totalItems, err := service.SearchByNameAndFilterByRating(page, perPage, name, ratingParam)

		assert.NoError(t, err)
		assert.Equal(t, foundProducts, result)
		assert.Equal(t, int64(len(foundProducts)), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Error case - Repository error", func(t *testing.T) {
		expectedLowerBound, expectedUpperBound, _ := getRatingBounds(ratingParam)

		repo.On("SearchByNameAndFilterByRating", page, perPage, name, ratingParam, expectedLowerBound, expectedUpperBound).Return(nil, int64(0), errors.New("some repository error")).Once()

		result, totalItems, err := service.SearchByNameAndFilterByRating(page, perPage, name, ratingParam)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Error case - Invalid Rating Param", func(t *testing.T) {
		ratingParam := "invalid"

		result, totalItems, err := service.SearchByNameAndFilterByRating(page, perPage, name, ratingParam)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
	})
}
