package service

import (
	"errors"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/mocks"
	userMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/mocks"
	user "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/service"
	utils "github.com/capstone-kelompok-7/backend-disappear/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestService(t *testing.T) (
	*AuthService,
	*mocks.RepositoryAuthInterface,
	*utils.JWTInterface,
	*userMocks.RepositoryUserInterface,
	*utils.HashInterface,
	*utils.CacheRepository) {

	repo := mocks.NewRepositoryAuthInterface(t)
	jwt := utils.NewJWTInterface(t)
	userRepo := userMocks.NewRepositoryUserInterface(t)
	hash := utils.NewHashInterface(t)
	cache := utils.NewCacheRepository(t)
	userService := user.NewUserService(userRepo, hash)
	service := NewAuthService(repo, jwt, userService, hash, cache)

	return service.(*AuthService), repo, jwt, userRepo, hash, cache
}

func TestAuthService_RegisterSocial(t *testing.T) {
	userID := uint64(1)
	request := &dto.RegisterSocialRequest{
		SocialID:     "123123131231",
		Provider:     "Google",
		Email:        "email@email.com",
		Name:         "Joni",
		PhotoProfile: "photo.jpg",
	}
	existingUser := &entities.UserModels{
		ID:           userID,
		SocialID:     "123123131231",
		Email:        "email@email.com",
		Name:         "Joni",
		PhotoProfile: "photo.jpg",
	}

	t.Run("Success Case - Register Social", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		authRepo.On("FindUserBySocialID", request.SocialID).Return(nil, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(existingUser, nil)

		result, err := authService.RegisterSocial(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Email Already Exists", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)

		result, err := authService.RegisterSocial(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "email sudah terdaftar")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Google ID Already Exists", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		authRepo.On("FindUserBySocialID", request.SocialID).Return(existingUser, nil)

		result, err := authService.RegisterSocial(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "Social ID sudah terdaftar")

		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Register Social Error", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _ := setupTestService(t)
		expectedError := errors.New("gagal mendaftarkan pengguna baru")
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		authRepo.On("FindUserBySocialID", request.SocialID).Return(nil, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(nil, expectedError)

		result, err := authService.RegisterSocial(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "gagal mendaftarkan pengguna baru")

		authRepo.AssertExpectations(t)
	})
}

func TestAuthService_LoginSocial(t *testing.T) {

	socialID := "someSocialID"

	user := &entities.UserModels{
		ID:        1,
		SocialID:  socialID,
		Email:     "test@example.com",
		Role:      "customer",
		LastLogin: time.Now(),
	}

	t.Run("Success Case - User Found and Login", func(t *testing.T) {
		authService, authRepo, jwtService, _, _, _ := setupTestService(t)
		authRepo.On("FindUserBySocialID", socialID).Return(user, nil)
		authRepo.On("UpdateLastLogin", user.ID, mock.AnythingOfType("time.Time")).Return(nil)
		jwtService.On("GenerateJWT", user.ID, user.Email, user.Role).Return("someAccessToken", nil)

		foundUser, accessToken, err := authService.LoginSocial(socialID)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "someAccessToken", accessToken)

		authRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		authService, authRepo, jwtService, _, _, _ := setupTestService(t)
		expectedError := errors.New("pengguna tidak ditemukan")
		authRepo.On("FindUserBySocialID", socialID).Return(nil, expectedError)

		foundUser, accessToken, err := authService.LoginSocial(socialID)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.EqualError(t, err, "pengguna tidak ditemukan")
		assert.Empty(t, accessToken)

		authRepo.AssertExpectations(t)
		jwtService.AssertNotCalled(t, "GenerateJWT")
	})

	// t.Run("Failed Case - Error Finding User by Social ID", func(t *testing.T) {
	// 	authService, authRepo, _, _, _, _ := setupTestService(t)
	// 	expectedErr := errors.New("some error here")
	// 	authRepo.On("FindUserBySocialID", socialID).Return(nil, expectedErr)

	// 	foundUser, _, err := authService.LoginSocial(socialID)

	// 	assert.Error(t, err)
	// 	assert.Nil(t, foundUser)
	// 	assert.Equal(t, err, expectedErr)

	// 	authRepo.AssertExpectations(t)
	// })

	t.Run("Failed Case - User is Nil After Finding", func(t *testing.T) {
		authService, authRepo, _, _, _, _ := setupTestService(t)
		authRepo.On("FindUserBySocialID", socialID).Return(nil, nil)

		foundUser, _, err := authService.LoginSocial(socialID)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, "pengguna tidak ditemukan", err.Error())

		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Updating LastLogin", func(t *testing.T) {
		authService, authRepo, jwtService, _, _, _ := setupTestService(t)
		authRepo.On("FindUserBySocialID", socialID).Return(user, nil)
		authRepo.On("UpdateLastLogin", user.ID, mock.AnythingOfType("time.Time")).Return(errors.New("update failed"))

		foundUser, accessToken, err := authService.LoginSocial(socialID)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.EqualError(t, err, "gagal memperbarui LastLogin")
		assert.Empty(t, accessToken)

		authRepo.AssertExpectations(t)
		jwtService.AssertNotCalled(t, "GenerateJWT")
	})

	t.Run("Failed Case - Error Generating JWT", func(t *testing.T) {
		authService, authRepo, jwtService, _, _, _ := setupTestService(t)
		user := &entities.UserModels{
			ID:        1,
			SocialID:  socialID,
			Email:     "test@example.com",
			Role:      "customer",
			LastLogin: time.Now(),
		}

		authRepo.On("FindUserBySocialID", socialID).Return(user, nil)
		authRepo.On("UpdateLastLogin", user.ID, mock.AnythingOfType("time.Time")).Return(nil)
		expectedErr := errors.New("JWT generation failed")
		jwtService.On("GenerateJWT", user.ID, user.Email, user.Role).Return("", expectedErr)

		foundUser, _, err := authService.LoginSocial(socialID)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, expectedErr, err)

		authRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}
