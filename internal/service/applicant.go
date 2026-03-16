package service

import "github.com/niiilov/go-dog-trapping/internal/dto"

func (s *Service) GetApplicantByDistrictID(districtID string) ([]*dto.Applicant, error) {
	applicants, err := s.repository.GetApplicantByDistrictID(districtID)
	if err != nil {
		return nil, err
	}
	return applicants, nil
}

func (s *Service) CreateApplicant(applicant *dto.CreateApplicantDTO) error {
	return s.repository.CreateApplicant(applicant)
}
