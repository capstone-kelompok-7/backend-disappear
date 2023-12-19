package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/caching"
	"github.com/labstack/gommon/log"

	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/email"
	"github.com/capstone-kelompok-7/backend-disappear/utils/otp"
)

type AuthService struct {
	repo        auth.RepositoryAuthInterface
	userService users.ServiceUserInterface
	jwt         utils.JWTInterface
	hash        utils.HashInterface
	cache       caching.CacheRepository
	email       email.EmailSenderInterface
}

func NewAuthService(repo auth.RepositoryAuthInterface, jwt utils.JWTInterface, userService users.ServiceUserInterface, hash utils.HashInterface, cache caching.CacheRepository, email email.EmailSenderInterface) auth.ServiceAuthInterface {
	return &AuthService{
		repo:        repo,
		jwt:         jwt,
		userService: userService,
		hash:        hash,
		cache:       cache,
		email:       email,
	}
}

func generateCacheKey(email, action string) string {
	return fmt.Sprintf("auth:%s:%s", email, action)
}

func (s *AuthService) Register(newData *entities.UserModels) (*entities.UserModels, error) {
	existingUser, _ := s.userService.GetUsersByEmail(newData.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashPassword, err := s.hash.GenerateHash(newData.Password)
	if err != nil {
		return nil, err
	}
	value := &entities.UserModels{
		Email:     newData.Email,
		Password:  hashPassword,
		Role:      "customer",
		Level:     "bronze",
		LastLogin: time.Now(),
	}

	result, err := s.repo.Register(value)
	if err != nil {
		return nil, err
	}

	generateOTP := otp.GenerateRandomOTP(6)
	newOTP := &entities.OTPModels{
		UserID:     int(result.ID),
		OTP:        generateOTP,
		ExpiredOTP: time.Now().Add(2 * time.Minute).Unix(),
	}

	_, errOtp := s.repo.SaveOTP(newOTP)
	if errOtp != nil {
		return nil, errOtp
	}
	err = s.email.EmailService(result.Email, generateOTP)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AuthService) Login(email, password, deviceToken string) (*entities.UserModels, string, error) {
	cachedToken, err := s.cache.Get(email)
	if err == nil {
		return nil, string(cachedToken), nil
	}
	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return nil, "", errors.New("user tidak ditemukan")
	}
	if !user.IsVerified {
		return nil, "", errors.New("akun anda belum diverifikasi")
	}
	isValidPassword, err := s.hash.ComparePassword(user.Password, password)
	if err != nil || !isValidPassword {
		return nil, "", errors.New("password salah")
	}

	user.LastLogin = time.Now()
	if err := s.repo.UpdateLastLogin(user.ID, user.LastLogin); err != nil {
		return nil, "", errors.New("gagal memperbarui LastLogin")
	}

	accessToken, err := s.jwt.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	cekDeviceToken, err := s.repo.CekDeviceTokenByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if cekDeviceToken != deviceToken {
		_, err := s.repo.UpdateDeviceTokenByID(email, deviceToken)
		if err != nil {
			return nil, "", errors.New("gagal memperbarui device token")
		}

	}

	err = s.cache.Set(email, []byte(accessToken), 1*time.Second)
	if err != nil {
		log.Error("Failed to store accessToken in cache:", err)
		return nil, "", errors.New("internal server error")
	}

	return user, accessToken, nil
}

func (s *AuthService) VerifyEmail(email, otp string) error {
	cacheKey := generateCacheKey(email, "verify_status")
	isVerified, err := s.cache.Get(cacheKey)
	if err == nil && string(isVerified) == "true" {
		return errors.New("email sudah diverifikasi sebelumnya")
	}

	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("user tidak ditemukan")
	}

	isValidOTP, err := s.repo.FindValidOTP(int(user.ID), otp)
	if err != nil {
		return err
	}

	if isValidOTP.ID == 0 {
		return errors.New("invalid atau OTP telah kadaluarsa")
	}

	user.IsVerified = true

	_, errUpdate := s.repo.UpdateUser(user)
	if errUpdate != nil {
		return errors.New("gagal verifikasi email")
	}

	errDeleteOTP := s.repo.DeleteOTP(isValidOTP)
	if errDeleteOTP != nil {
		return errors.New("gagal delete OTP")
	}

	err = s.cache.Set(cacheKey, []byte("true"), 1*time.Second)
	if err != nil {
		return errors.New("gagal menyimpan status verifikasi email ke cache")
	}

	return nil
}

func (s *AuthService) ResendOTP(email string) (*entities.OTPModels, error) {
	cacheKey := generateCacheKey(email, "verify_status")
	isVerified, err := s.cache.Get(cacheKey)
	if err == nil && string(isVerified) == "true" {
		return nil, errors.New("email sudah diverifikasi, tidak dapat mengirim ulang OTP")
	}

	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	errDeleteOTP := s.repo.DeleteUserOTP(user.ID)
	if errDeleteOTP != nil {
		return nil, errDeleteOTP
	}
	generateOTP := otp.GenerateRandomOTP(6)
	newOTP := &entities.OTPModels{
		UserID:     int(user.ID),
		OTP:        generateOTP,
		ExpiredOTP: time.Now().Add(2 * time.Minute).Unix(),
	}

	_, err = s.repo.SaveOTP(newOTP)
	if err != nil {
		return nil, err
	}
	return newOTP, nil
}

func (s *AuthService) ResetPassword(email, newPassword, confirmPass string) error {
	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	if user.ID == 0 {
		return errors.New("user tidak ditemukan")
	}

	if newPassword != confirmPass {
		return errors.New("konfirmasi password tidak cocok")
	}

	newPasswordHash, err := s.hash.GenerateHash(newPassword)
	if err != nil {
		return errors.New("gagal melakukan hash password baru")
	}

	err = s.repo.ResetPassword(email, newPasswordHash)
	if err != nil {
		return errors.New("gagal reset pass: ")
	}
	return nil
}

func (s *AuthService) VerifyOTP(email, otp string) (string, error) {
	emailVerifyCacheKey := generateCacheKey(email, "verify_status")
	isVerified, err := s.cache.Get(emailVerifyCacheKey)
	if err == nil && string(isVerified) == "true" {
		return "", errors.New("email sudah diverifikasi")
	}

	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user tidak ditemukan")
	}

	accessTokenCacheKey := generateCacheKey(email, "access_token")
	cachedToken, err := s.cache.Get(accessTokenCacheKey)
	if err == nil {
		return string(cachedToken), nil
	}

	isValidOTP, err := s.repo.FindValidOTP(int(user.ID), otp)
	if err != nil {
		return "", err
	}

	if isValidOTP.ID == 0 {
		return "", errors.New("invalid atau OTP telah kadaluarsa")
	}

	user.IsVerified = true

	_, errUpdate := s.repo.UpdateUser(user)
	if errUpdate != nil {
		return "", errors.New("gagal verifikasi email")
	}

	errDeleteOTP := s.repo.DeleteOTP(isValidOTP)
	if errDeleteOTP != nil {
		return "", errors.New("gagal delete OTP")
	}

	accessToken, err := s.jwt.GenerateJWT(user.ID, user.Role, user.Email)
	if err != nil {
		return "", errors.New("gagal generate access token")
	}

	err = s.cache.Set(accessTokenCacheKey, []byte(accessToken), 1*time.Second)
	if err != nil {
		return "", errors.New("gagal menyimpan access token ke cache")
	}

	err = s.cache.Set(emailVerifyCacheKey, []byte("true"), 1*time.Second)
	if err != nil {
		return "", errors.New("gagal menyimpan status verifikasi email ke cache")
	}

	return accessToken, nil
}

func (s *AuthService) RegisterSocial(req *dto.RegisterSocialRequest) (*entities.UserModels, error) {
	existingUser, _ := s.userService.GetUsersByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	existingUserBySocialID, _ := s.repo.FindUserBySocialID(req.SocialID)
	if existingUserBySocialID != nil {
		return nil, errors.New("Social ID sudah terdaftar")
	}

	value := &entities.UserModels{
		SocialID:     req.SocialID,
		Provider:     req.Provider,
		Email:        req.Email,
		Name:         req.Name,
		PhotoProfile: req.PhotoProfile,
		Role:         "customer",
		Level:        "bronze",
		LastLogin:    time.Time{},
	}

	result, err := s.repo.Register(value)
	if err != nil {
		return nil, errors.New("gagal mendaftarkan pengguna baru")
	}
	return result, nil
}

func (s *AuthService) LoginSocial(socialID string) (*entities.UserModels, string, error) {
	user, err := s.repo.FindUserBySocialID(socialID)
	if err != nil {
		return nil, "", errors.New("pengguna tidak ditemukan")
	}

	if user == nil {
		return nil, "", errors.New("pengguna tidak ditemukan")
	}

	user.LastLogin = time.Now()
	if err := s.repo.UpdateLastLogin(user.ID, user.LastLogin); err != nil {
		return nil, "", errors.New("gagal memperbarui LastLogin")
	}

	accessToken, err := s.jwt.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}
