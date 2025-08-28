package service

import (
	"errors"

	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/pkg/security"
)

//тут доп обработка по хуйне

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)

type repository interface {
	CreateAccount(account *dto.Account) (string, error)
	ValidateAccount(account *dto.Account) (id string, hashPass string, err error)
}
type Service struct {
	repository repository
}

func New(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateAccount(account *dto.Account) (string, error) {
	hashPass, err := security.Encode(account.Password)
	if err != nil {
		return "", err
	}
	account.Password = hashPass

	return s.repository.CreateAccount(account)
}
func (s *Service) ValidateAccount(account *dto.Account) (id string, err error) {
	id, hashPass, err := s.repository.ValidateAccount(account)
	if err != nil {
		return "", err
	}

	if !security.Check(account.Password, hashPass) {
		return "", ErrInvalidPassword
	}
	return id, nil
}
