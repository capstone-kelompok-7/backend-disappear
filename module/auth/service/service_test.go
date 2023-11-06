package service

import (
	"errors"
	authMock "github.com/capstone-kelompok-7/backend-disappear/module/auth/mocks"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	userMock "github.com/capstone-kelompok-7/backend-disappear/module/users/mocks"
	service2 "github.com/capstone-kelompok-7/backend-disappear/module/users/service"
	utilsMock "github.com/capstone-kelompok-7/backend-disappear/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegister(t *testing.T) {
	jwt := utilsMock.NewJWTInterface(t)
	repo := authMock.NewRepositoryAuthInterface(t)
	userRepo := userMock.NewRepositoryUserInterface(t)
	hash := utilsMock.NewHashInterface(t)
	userService := service2.NewUserService(userRepo)
	service := NewAuthService(repo, jwt, userService, hash)
	newUser := domain.UserModels{
		Email:    "user@mail.com",
		Phone:    "08123123123123",
		Password: "a",
	}
	t.Run("Kasus Hash Password Gagal", func(t *testing.T) {
		expectedHashedPassword := "asdjsdhasdasdasj"
		hash.On("GenerateHash", newUser.Password).Return(expectedHashedPassword, errors.New("kesalahan pada hash password")).Once()

		result, err := service.Register(&newUser)
		assert.Error(t, err)
		assert.Nil(t, result)

		repo.AssertExpectations(t)

		repo.AssertNotCalled(t, "Register", mock.AnythingOfType("*domain.UserModels"))
	})

	t.Run("Kasus Sukses", func(t *testing.T) {
		expectedPassword := "hashed_password"
		hash.On("GenerateHash", newUser.Password).Return(expectedPassword, nil).Once()
		repo.On("Register", mock.AnythingOfType("*domain.UserModels")).Return(&domain.UserModels{
			Email:    newUser.Email,
			Phone:    newUser.Phone,
			Password: expectedPassword,
			Role:     "customer",
		}, nil).Once()

		result, err := service.Register(&newUser)
		assert.Nil(t, err)
		assert.Equal(t, newUser.Phone, result.Phone)
		assert.Equal(t, newUser.Email, result.Email)
		assert.Equal(t, expectedPassword, result.Password)
		assert.Equal(t, "customer", result.Role)

		hash.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Kasus Gagal", func(t *testing.T) {
		expectedErr := errors.New("registration failed")
		expectedHashedPassword := "asdjsdhasdasdasj"
		hash.On("GenerateHash", newUser.Password).Return(expectedHashedPassword, nil).Once()
		repo.On("Register", mock.AnythingOfType("*domain.UserModels")).Return(nil, expectedErr).Once()

		result, err := service.Register(&newUser)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		hash.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	jwt := utilsMock.NewJWTInterface(t)
	repo := authMock.NewRepositoryAuthInterface(t)
	userRepo := userMock.NewRepositoryUserInterface(t)
	hash := utilsMock.NewHashInterface(t)
	userService := service2.NewUserService(userRepo)
	service := NewAuthService(repo, jwt, userService, hash)
	email := "user@example.com"
	password := "hashed_password"

	existingUser := &domain.UserModels{
		ID:       1,
		Email:    "user@example.com",
		Password: "hashed_password",
		Role:     "customer",
	}

	t.Run("Kasus Sukses Login", func(t *testing.T) {
		expectedAccessToken := "mocked-access-token"

		userRepo.On("GetUsersByEmail", email).Return(existingUser, nil).Once()
		hash.On("ComparePassword", existingUser.Password, password).Return(true, nil).Once()
		jwt.On("GenerateJWT", existingUser.ID, existingUser.Role).Return(expectedAccessToken, nil).Once()

		resultUser, accessToken, err := service.Login(email, password)
		assert.Nil(t, err)
		assert.Equal(t, existingUser, resultUser)
		assert.Equal(t, expectedAccessToken, accessToken)

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("Kasus Kesalahan Password", func(t *testing.T) {

		userRepo.On("GetUsersByEmail", email).Return(existingUser, nil).Once()
		hash.On("ComparePassword", existingUser.Password, password).Return(false, nil).Once()

		resultUser, accessToken, err := service.Login(email, password)
		assert.Error(t, err)
		assert.Nil(t, resultUser)
		assert.Empty(t, accessToken)

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("Kasus Kesalahan Pengguna Tidak Ditemukan", func(t *testing.T) {

		userRepo.On("GetUsersByEmail", email).Return(nil, errors.New("pengguna tidak ditemukan")).Once()

		resultUser, accessToken, err := service.Login(email, password)

		assert.Error(t, err)
		assert.Nil(t, resultUser)
		assert.Empty(t, accessToken)

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("Kasus Kesalahan Gagal Generate Token", func(t *testing.T) {

		jwtError := errors.New("gagal generate jwt")
		userRepo.On("GetUsersByEmail", email).Return(existingUser, nil).Once()
		hash.On("ComparePassword", existingUser.Password, password).Return(true, nil).Once()
		jwt.On("GenerateJWT", existingUser.ID, existingUser.Role).Return("", jwtError).Once()

		resultUser, accessToken, err := service.Login(email, password)
		assert.Error(t, err)
		assert.Nil(t, resultUser)
		assert.Empty(t, accessToken)

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

}
