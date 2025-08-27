package service

import "github.com/niiilov/go-dog-trapping/internal/dto"

//тут доп обработка по хуйне

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) CreateAccount(account *dto.Account, uuid string) error {
	return nil
}
func (s *Service) ValidateAccount(account *dto.Account) (id string, err error) {
	return "", nil
}
