package service

import (
	"github.com/niiilov/go-dog-trapping/internal/dto"
	validate "github.com/niiilov/go-dog-trapping/pkg/validator"
)

func (s *Service) CreateTerOtdel(terOtdel *dto.CreateTerOtdelDTO) (string, error) {
	err := validate.Validate(terOtdel)
	if err != nil {
		return "", ErrInvalidData
	}
	return s.repository.CreateTerOtdel(terOtdel)
}

func (s *Service) GetTerOtdels() ([]*dto.TerOtdel, error) {
	return s.repository.GetTerOtdels()
}

func (s *Service) GetTerrOtdelsByDistrictID(district_id string) ([]*dto.TerOtdel, error) {
	return s.repository.GetTerrOtdelsByDistrictID(district_id)
}

func (s *Service) DeleteTerOtdel(id string) error {
	return s.repository.DeleteTerOtdel(id)
}
