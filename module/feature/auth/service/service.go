package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	"time"

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
}

func NewAuthService(repo auth.RepositoryAuthInterface, jwt utils.JWTInterface, userService users.ServiceUserInterface, hash utils.HashInterface) auth.ServiceAuthInterface {
	return &AuthService{
		repo:        repo,
		jwt:         jwt,
		userService: userService,
		hash:        hash,
	}
}

func (s *AuthService) Register(newData *entities.UserModels) (*entities.UserModels, error) {
	hashPassword, err := s.hash.GenerateHash(newData.Password)
	if err != nil {
		return nil, err
	}
	value := &entities.UserModels{
		Email:    newData.Email,
		Password: hashPassword,
		Role:     "customer",
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
	err = email.EmaiilService(result.Email, generateOTP)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AuthService) Login(email, password string) (*entities.UserModels, string, error) {
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

func (s *AuthService) ResendOTP(email string) (*entities.OTPModels, error) {
	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user tidak ditemukan pada email ini")
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
		return err
	}
	if user.ID == 0 {
		return errors.New("user tidak ditemukan")
	}

	if newPassword != confirmPass {
		return errors.New("konfirmasi password tidak cocok")
	}

	newPasswordHash, err := s.hash.GenerateHash(newPassword)
	if err != nil {
		return fmt.Errorf("gagal melakukan hash password baru: %w", err)
	}

	err = s.repo.ResetPassword(email, newPasswordHash)
	if err != nil {
		return fmt.Errorf("gagal reset pass: %w", err)
	}
	return nil
}
