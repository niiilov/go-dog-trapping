package service

import (
	"errors"
	"strconv"

	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/internal/worker"
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
	ValidateAccount(account *dto.AuthCredentials) (id string, hashPass string, err error)
	SendRequest(request *dto.RequestFull) (int, error)
	GetAllRequests() ([]*dto.RequestFull, error)
	GetRequestsByOtdel(otdel_id string) ([]*dto.RequestFull, error)
	ChangePassword(req *dto.ChangePasswordRequest) error
	GetUserProfile(userId string) (*dto.UserProfile, error)
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
func (s *Service) ValidateAccount(account *dto.AuthCredentials) (string, error) {

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

	number, err := s.repository.SendRequest(request)
	if err != nil {
		return err
	}

	if number == 0 {
		return ErrInvalidData
	}
	reqData := &dto.RequestForGenerating{
		Number:        strconv.Itoa(number),
		Applicant:     request.Applicant.Name,
		Source:        request.Source.Name,
		Address:       request.Address,
		DogsCount:     request.DogsCount,
		Behavior:      request.Behavior,
		Urgency:       request.Urgency,
		ContactPerson: request.ContactPerson,
	}

	if err = worker.SendRequestForGenerating(reqData); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRequestsByOtdel(otdel_id string) ([]*dto.RequestFull, error) {
	requests, err := s.repository.GetRequestsByOtdel(otdel_id)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *Service) GetAllRequests() ([]*dto.RequestFull, error) {
	requests, err := s.repository.GetAllRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *Service) ChangePassword(req *dto.ChangePasswordRequest) error {
	var valid dto.AuthCredentials

	valid.Login = req.Login
	valid.Password = req.OldPassword

	_, hashPass, err := s.repository.ValidateAccount(&valid)
	if err != nil {
		return ErrInvalidData
	}
	if !security.Check(valid.Password, hashPass) {
		return ErrInvalidPassword
	}

	req.NewPassword, err = security.Encode(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.repository.ChangePassword(req)

	return err

}

func (s *Service) GetUserProfile(userId string) (*dto.UserProfile, error) {
	profile, err := s.repository.GetUserProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
