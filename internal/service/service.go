package service

import (
	"errors"

	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/pkg/security"
	validate "github.com/niiilov/go-dog-trapping/pkg/validator"
)

//тут доп обработка по хуйне

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidData     = errors.New("invalid data")
)

type repository interface {
	CreateAccount(account *dto.Account) (string, error)
	ValidateAccount(account *dto.Account) (id string, hashPass string, err error)
	SendRequest(request *dto.RequestFull) error
	GetRequests() ([]*dto.RequestFull, error)
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
	err := validate.Validate(account)
	if err != nil {

		return "", ErrInvalidData
	}

	hashPass, err := security.Encode(account.Password)
	if err != nil {
		return "", err
	}
	account.Password = hashPass

	return s.repository.CreateAccount(account)
}
func (s *Service) ValidateAccount(account *dto.Account) (string, error) {

	err := validate.Validate(account)
	if err != nil {

		return "", ErrInvalidData
	}

	id, hashPass, err := s.repository.ValidateAccount(account)
	if err != nil {
		return "", err
	}

	if !security.Check(account.Password, hashPass) {
		return "", ErrInvalidPassword
	}
	return id, nil
}

func (s *Service) SendRequest(request *dto.RequestFull) error {
	err := validate.Validate(request)
	if err != nil {
		return ErrInvalidData
	}

	if err := s.repository.SendRequest(request); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRequests() ([]*dto.RequestFull, error) {
	requests, err := s.repository.GetRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}
