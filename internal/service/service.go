package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

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
	ValidateAccount(account *dto.AuthCredentials) (user *dto.UserProfile, hashPass string, role_id string, err error)
	SendRequest(request *dto.RequestFull) (int, error)
	GetAllRequests(year string) ([]*dto.RequestFull, error)
	GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error)
	ChangePassword(req *dto.ChangePasswordRequest) error
	ChangeProfileInfo(req *dto.ChangeProfileRequest) error
	GetUserProfile(userId string) (*dto.UserProfile, error)
	ChangeStatusRequest(req *dto.ChangeStatusRequest) error

	GetRequestsByNumber(request_number []string) ([]*dto.RequestFull, error)
}

type storage interface {
	UploadFile(ctx context.Context, objectKey string, fileName string) error
	GetFileURL(objectKey string) string
}
type Service struct {
	repository repository
	storage    storage
}

func New(r repository, s storage) *Service {
	return &Service{
		repository: r,
		storage:    s,
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
func (s *Service) ValidateAccount(account *dto.AuthCredentials) (*dto.UserProfile, string, error) {

	err := validate.Validate(account)
	if err != nil {

		return nil, "", ErrInvalidData
	}

	user, hashPass, role_id, err := s.repository.ValidateAccount(account)
	if err != nil {
		return nil, "", err
	}

	if !security.Check(account.Password, hashPass) {
		return nil, "", ErrInvalidPassword
	}
	return user, role_id, nil
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
	sharedDir := "/app/shared/"
	key := "zayavka_" + strconv.Itoa(number) + "_" + time.Now().Format("2006") + ".xlsx"
	filename := sharedDir + key

	if err = s.storage.UploadFile(context.TODO(), key, filename); err != nil {
		//лог
		fmt.Println("Error upload file to S3:", err)
		return err
	}

	os.Remove(filename)

	return nil
}

func (s *Service) ChangeStatusRequest(req *dto.ChangeStatusRequest) error {

	err := s.repository.ChangeStatusRequest(req)
	return err
}

func (s *Service) GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error) {
	requests, err := s.repository.GetRequestsByOtdel(otdel_id, year)
	if err != nil {
		return nil, err
	}
	go s.validateDelay(requests)
	return requests, nil
}

func (s *Service) GetAllRequests(year string) ([]*dto.RequestFull, error) {
	requests, err := s.repository.GetAllRequests(year)
	if err != nil {
		return nil, err
	}
	go s.validateDelay(requests)
	return requests, nil
}

func (s *Service) ChangeProfileInfo(req *dto.ChangeProfileRequest) error {
	err := s.repository.ChangeProfileInfo(req)
	return err
}

func (s *Service) ChangePassword(req *dto.ChangePasswordRequest) error {
	var valid dto.AuthCredentials

	valid.Login = req.Login
	valid.Password = req.OldPassword

	_, hashPass, _, err := s.repository.ValidateAccount(&valid)
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

func (s *Service) GetFileURL(objectKey string) string {
	return s.storage.GetFileURL(objectKey)
}

func (s *Service) UploadAct(req *dto.UploadActRequests, key string, filename string) error {
	if err := s.storage.UploadFile(context.TODO(), key, filename); err != nil {
		//лог
		fmt.Println("Error upload file to S3:", err)
		return err
	}
	// Удаляем временный файл
	os.Remove(filename)
	var statusReq dto.ChangeStatusRequest
	statusReq.ID = req.ID
	statusReq.Status = req.Status
	if err := s.ChangeStatusRequest(&statusReq); err != nil {
		return err
	}

	return nil
}

func (s *Service) validateDelay(requests []*dto.RequestFull) error {
	twoWeeks := 14 * 24 * time.Hour

	for _, req := range requests {
		if req == nil {
			continue
		}
		if req.Status == "Завершена" || req.Status == "Просрочена" {
			continue
		}
		// если прошло больше двух недель с момента создания — меняем статус
		if time.Since(req.CreatedAt) > twoWeeks {

			statusReq := &dto.ChangeStatusRequest{
				ID:     req.ID,
				Status: "Просрочена", // <-- поменяйте на нужный вам статус
			}

			if err := s.repository.ChangeStatusRequest(statusReq); err != nil {
				fmt.Println("failed to change status for request", req.ID, ":", err)
				// продолжаем обработку остальных заявок
			}
		}
	}

	return nil

}

func (s *Service) GenerateMultiple(req *dto.GenerateMultipleRequest) (string, error) {
	requests, err := s.repository.GetRequestsByNumber(req.Numbers)
	if err != nil {
		return "", err
	}

	reqs := dto.RequestForGeneratingSomething{
		Requests: requests,
		Number:   requests[0].Number,
	}
	if err := worker.SendRequestForGeneratingSomething(&reqs); err != nil {
		return "", err
	}
	number := requests[0].Number
	sharedDir := "/app/shared/"
	key := "zayavka_" + number + "_" + time.Now().Format("2006") + ".xlsx"
	filename := sharedDir + key

	if err = s.storage.UploadFile(context.TODO(), key, filename); err != nil {
		//лог
		fmt.Println("Error upload file to S3:", err)
		return "", err
	}

	os.Remove(filename)

	url := s.storage.GetFileURL(filename)

	return url, nil
}
