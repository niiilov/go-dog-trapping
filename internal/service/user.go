package service

import (
	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/pkg/security"
	validate "github.com/niiilov/go-dog-trapping/pkg/validator"
)

func (r *Service) CreateUser(user *dto.CreateUserDTO) error {

	passwordHash, err := security.Encode(user.Password)
	if err != nil {
		return err
	}

	err = r.repository.CreateUser(user, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUsers() ([]*dto.GetUserDTO, error) {
	users, err := s.repository.GetUsers()
	if err != nil {
		return nil, err
	}

	var result []*dto.GetUserDTO
	for _, user := range users {
		result = append(result, &dto.GetUserDTO{
			ID:         user.ID,
			FullName:   user.FullName,
			Login:      user.Login,
			RoleID:     user.RoleID,
			DistrictID: user.DistrictID,
			TerOtdelID: user.TerOtdelID,
		})
	}
	return result, nil
}

func (s *Service) DeleteUser(userID string) error {
	return s.repository.DeleteUser(userID)
}

func (s *Service) ValidateAccount(account *dto.AuthCredentials) (*dto.UserProfile, error) {

	err := validate.Validate(account)
	if err != nil {

		return nil, ErrInvalidData
	}

	user, err := s.repository.GetUserByLogin(account.Login)
	if err != nil {
		return nil, err
	}

	if !security.Check(account.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	profile, err := s.repository.GetUserProfile(user.ID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
