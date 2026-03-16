package service

import (
	"github.com/niiilov/go-dog-trapping/internal/dto"
	validate "github.com/niiilov/go-dog-trapping/pkg/validator"
)

func (s *Service) CreateDistrict(district *dto.District) error {
	err := validate.Validate(district)
	if err != nil {
		return ErrInvalidData
	}
	return s.repository.CreateDistrict(district)
}

func (s *Service) GetDistricts() ([]*dto.District, error) {
	return s.repository.GetDistricts()
}

func (s *Service) DeleteDistrict(id string) error {
	return s.repository.DeleteDistrict(id)
}
