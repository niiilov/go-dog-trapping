package service

import "github.com/niiilov/go-dog-trapping/internal/dto"

func (s *Service) GetRoles() ([]*dto.Role, error) {
	roles, err := s.repository.GetRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
