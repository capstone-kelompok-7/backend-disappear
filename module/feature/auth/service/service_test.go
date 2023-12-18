package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/utils/otp"
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
	*utils.CacheRepository,
	*utils.EmailSenderInterface) {

	repo := mocks.NewRepositoryAuthInterface(t)
	jwt := utils.NewJWTInterface(t)
	userRepo := userMocks.NewRepositoryUserInterface(t)
	hash := utils.NewHashInterface(t)
	cache := utils.NewCacheRepository(t)
	userService := user.NewUserService(userRepo, hash)
	email := utils.NewEmailSenderInterface(t)
	service := NewAuthService(repo, jwt, userService, hash, cache, email)

	return service.(*AuthService), repo, jwt, userRepo, hash, cache, email
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
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
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
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)

		result, err := authService.RegisterSocial(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "email sudah terdaftar")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Google ID Already Exists", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		authRepo.On("FindUserBySocialID", request.SocialID).Return(existingUser, nil)

		result, err := authService.RegisterSocial(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "Social ID sudah terdaftar")

		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Register Social Error", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
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
		authService, authRepo, jwtService, _, _, _, _ := setupTestService(t)
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
		authService, authRepo, jwtService, _, _, _, _ := setupTestService(t)
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

	t.Run("Failed Case - User is Nil After Finding", func(t *testing.T) {
		authService, authRepo, _, _, _, _, _ := setupTestService(t)
		authRepo.On("FindUserBySocialID", socialID).Return(nil, nil)

		foundUser, _, err := authService.LoginSocial(socialID)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, "pengguna tidak ditemukan", err.Error())

		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Updating LastLogin", func(t *testing.T) {
		authService, authRepo, jwtService, _, _, _, _ := setupTestService(t)
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
		authService, authRepo, jwtService, _, _, _, _ := setupTestService(t)

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

func TestGenerateCacheKey(t *testing.T) {
	email := "test@example.com"
	action := "reset_password"
	expectedKey := "auth:test@example.com:reset_password"

	result := generateCacheKey(email, action)

	assert.Equal(t, expectedKey, result, "Generated cache key is incorrect")
}

func TestAuthService_Register(t *testing.T) {

	userID := uint64(1)
	passwordHash := "hashedpassword"
	request := &entities.UserModels{
		Email:    "email@email.com",
		Password: "admin123",
	}
	existingUser := &entities.UserModels{
		ID:       userID,
		Email:    "email@email.com",
		Password: "admin123",
	}
	expectedErr := errors.New("some error")

	t.Run("Failed Case - Email Already Exists", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)
		result, err := authService.Register(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "email sudah terdaftar")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Generate Hash", func(t *testing.T) {
		authService, authRepo, _, userRepo, hash, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		hash.On("GenerateHash", request.Password).Return("", expectedErr)
		result, err := authService.Register(request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, err, expectedErr)

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Register", func(t *testing.T) {
		authService, authRepo, _, userRepo, hash, _, email := setupTestService(t)

		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		hash.On("GenerateHash", request.Password).Return(passwordHash, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(nil, expectedErr)

		result, err := authService.Register(request)

		assert.Error(t, err)
		assert.Nil(t, result)

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		email.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Save OTP", func(t *testing.T) {
		authService, authRepo, _, userRepo, hash, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		hash.On("GenerateHash", request.Password).Return(passwordHash, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(&entities.UserModels{ID: userID}, nil)
		authRepo.On("SaveOTP", mock.AnythingOfType("*entities.OTPModels")).Return(nil, errors.New("failed to save OTP"))

		_, err := authService.Register(request)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to save OTP")

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Email Sending Error", func(t *testing.T) {
		authService, authRepo, _, userRepo, hash, _, emailService := setupTestService(t)

		generateOTP := otp.GenerateRandomOTP(6)
		newOTP := &entities.OTPModels{
			UserID: int(userID),
			OTP:    generateOTP,
		}

		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		hash.On("GenerateHash", request.Password).Return(passwordHash, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(&entities.UserModels{ID: userID}, nil)
		authRepo.On("SaveOTP", mock.AnythingOfType("*entities.OTPModels")).Return(newOTP, nil)
		emailService.On("EmailService", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(errors.New("failed to send email")).Once()

		_, err := authService.Register(request)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to send email")

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		emailService.AssertExpectations(t)
	})

	t.Run("Success Case - Register", func(t *testing.T) {
		authService, authRepo, _, userRepo, hash, _, emailService := setupTestService(t)

		generateOTP := otp.GenerateRandomOTP(6)
		newOTP := &entities.OTPModels{
			UserID: int(userID),
			OTP:    generateOTP,
		}

		userRepo.On("GetUsersByEmail", request.Email).Return(nil, nil)
		hash.On("GenerateHash", request.Password).Return(passwordHash, nil)
		authRepo.On("Register", mock.AnythingOfType("*entities.UserModels")).Return(&entities.UserModels{ID: userID}, nil)
		authRepo.On("SaveOTP", mock.AnythingOfType("*entities.OTPModels")).Return(newOTP, nil)
		emailService.On("EmailService", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil).Once()

		result, err := authService.Register(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		userRepo.AssertExpectations(t)
		hash.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		emailService.AssertExpectations(t)
	})

}

func TestAuthService_ResetPassword(t *testing.T) {
	authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)
	userID := uint64(1)
	request := &entities.UserModels{
		Email:    "email@email.com",
		Password: "admin123",
	}
	email := "email@email.com"
	password := "admin123"
	confirmPass := "admin123"
	existingUser := &entities.UserModels{
		ID:       userID,
		Email:    "email@email.com",
		Password: "admin123",
	}

	t.Run("Failed Case - Email Not Found", func(t *testing.T) {
		expectedErr := errors.New("user tidak ditemukan")
		userRepo.On("GetUsersByEmail", request.Email).Return(nil, expectedErr)

		err := authService.ResetPassword(email, password, confirmPass)

		assert.Error(t, err)
		assert.EqualError(t, err, "user tidak ditemukan")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - User ID Zero", func(t *testing.T) {
		user := &entities.UserModels{
			ID:       0,
			Email:    "email@email.com",
			Password: "hashedpassword",
		}
		userRepo.On("GetUsersByEmail", request.Email).Return(user, nil)

		err := authService.ResetPassword(email, password, confirmPass)

		assert.Error(t, err)
		assert.EqualError(t, err, "user tidak ditemukan")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Passwords Mismatch", func(t *testing.T) {
		authService, authRepo, _, userRepo, _, _, _ := setupTestService(t)

		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)

		err := authService.ResetPassword(email, "newpassword", confirmPass)

		assert.Error(t, err)
		assert.EqualError(t, err, "konfirmasi password tidak cocok")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Hash Password Error", func(t *testing.T) {
		authService, authRepo, _, userRepo, mockHash, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)
		mockHash.On("GenerateHash", mock.AnythingOfType("string")).Return("", errors.New("hashing error"))

		err := authService.ResetPassword(email, password, confirmPass)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal melakukan hash password baru")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		mockHash.AssertCalled(t, "GenerateHash", password)
	})

	t.Run("Failed Case - Reset Password Error", func(t *testing.T) {
		authService, authRepo, _, userRepo, mockHash, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)
		mockHash.On("GenerateHash", mock.AnythingOfType("string")).Return("hashedpassword", nil)
		authRepo.On("ResetPassword", email, "hashedpassword").Return(errors.New("reset error"))

		err := authService.ResetPassword(email, password, confirmPass)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal reset pass: ")

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		mockHash.AssertCalled(t, "GenerateHash", password)
	})

	t.Run("Success Case - Password Reset", func(t *testing.T) {
		authService, authRepo, _, userRepo, mockHash, _, _ := setupTestService(t)
		userRepo.On("GetUsersByEmail", request.Email).Return(existingUser, nil)
		mockHash.On("GenerateHash", mock.AnythingOfType("string")).Return("newhashedpassword", nil)
		authRepo.On("ResetPassword", email, "newhashedpassword").Return(nil)

		err := authService.ResetPassword(email, password, confirmPass)

		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
		authRepo.AssertExpectations(t)
		mockHash.AssertCalled(t, "GenerateHash", password)
	})

}
