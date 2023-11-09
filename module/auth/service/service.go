package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/email"
	"github.com/capstone-kelompok-7/backend-disappear/utils/otp"
	"time"
)

type AuthService struct {
	repo        auth.RepositoryAuthInterface
	userService users.ServiceUserInterface
	jwt         utils.JWTInterface
	hash        utils.HashInterface
}

func NewAuthService(repo auth.RepositoryAuthInterface, jwt utils.JWTInterface, userService users.ServiceUserInterface, hash utils.HashInterface) auth.ServiceAuthInterface {
	return &AuthService{
		repo:        repo,
		jwt:         jwt,
		userService: userService,
		hash:        hash,
	}
}

func (s *AuthService) Register(newData *domain.UserModels) (*domain.UserModels, error) {
	hashPassword, err := s.hash.GenerateHash(newData.Password)
	if err != nil {
		return nil, err
	}
	value := &domain.UserModels{
		Email:    newData.Email,
		Password: hashPassword,
		Role:     "customer",
	}

	result, err := s.repo.Register(value)
	if err != nil {
		return nil, err
	}

	generateOTP := otp.GenerateRandomOTP(6)
	newOTP := &domain.OTPModels{
		UserID:     int(result.ID),
		OTP:        generateOTP,
		ExpiredOTP: time.Now().Add(2 * time.Minute).Unix(),
	}

	_, errOtp := s.repo.SaveOTP(newOTP)
	if errOtp != nil {
		return nil, errOtp
	}
	err = email.EmaiilService(result.Email, generateOTP)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AuthService) Login(email, password string) (*domain.UserModels, string, error) {
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

	accessToken, err := s.jwt.GenerateJWT(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}

func (s *AuthService) VerifyEmail(email, otp string) error {
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
		return errors.New("Invalid atau OTP telah kadaluarsa")
	}

	user.IsVerified = true

	_, errUpdate := s.repo.UpdateUser(user)
	if errUpdate != nil {
		return errors.New("Gagal verifikasi email")
	}

	errDeleteOTP := s.repo.DeleteOTP(isValidOTP)
	if errDeleteOTP != nil {
		return errors.New("Gagal delete OTP")
	}

	return nil
}
