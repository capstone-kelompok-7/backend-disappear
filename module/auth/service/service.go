package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
)

type AuthService struct {
	repo        auth.RepositoryAuthInterface
	userService users.ServiceUserInterface
	utils       utils.JWTInterface
}

func NewAuthService(repo auth.RepositoryAuthInterface, utils utils.JWTInterface, userService users.ServiceUserInterface) auth.ServiceAuthInterface {
	return &AuthService{
		repo:        repo,
		utils:       utils,
		userService: userService,
	}
}

func (s *AuthService) Register(newData *domain.UserModels) (*domain.UserModels, error) {
	hashPassword, err := utils.GenerateHash(newData.Password)
	if err != nil {
		return nil, err
	}
	value := &domain.UserModels{
		Email:    newData.Email,
		Phone:    newData.Phone,
		Password: hashPassword,
		Role:     "customer",
	}

	result, err := s.repo.Register(value)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *AuthService) Login(email, password string) (*domain.UserModels, string, error) {
	user, err := s.userService.GetUsersByEmail(email)
	if err != nil {
		return nil, "", err
	}

	isValidPassword, err := utils.ComparePassword(user.Password, password)
	if err != nil || !isValidPassword {
		return nil, "", errors.New("invalid password")
	}

	accessToken, err := s.utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}
