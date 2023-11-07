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
