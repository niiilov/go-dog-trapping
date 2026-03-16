package service

import (
	"context"
	"errors"

	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/pkg/security"
)

//тут доп обработка по хуйне

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidData     = errors.New("invalid data")
)

type repository interface {
	// CreateAccount(account *dto.Account) (string, error)
	// ValidateAccount(account *dto.AuthCredentials) (user *dto.UserProfile, hashPass string, role_id string, sourceID string, err error)
	// SendRequest(request *dto.RequestFull) (int, error)
	// GetAllRequests(year string) ([]*dto.RequestFull, error)
	// GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error)
	// ChangePassword(req *dto.ChangePasswordRequest) error
	// ChangeProfileInfo(req *dto.ChangeProfileRequest) error

	// //ChangeStatusRequest(req *dto.ChangeStatusRequest) error
	// //DeleteRequest(reqID string) error
	// GetRequestsByDate(dateFrom *time.Time, dateTo *time.Time) ([]*dto.RequestForGenerating, error)
	// GetRequestsByIDs(reqIDs []string) ([]*dto.RequestForGenerating, error)

	// NewNotPemamentApplicant(name string) (string, error)

	CreateExternalUser(user *dto.ExternalUser) (int, error)
	GetExternalUserByID(id int) (*dto.ExternalUser, error)

	ChangeAcceptedStatusExternalUser(id int, isAccepted bool) error

	AcceptExternalUser(user *dto.AcceptExternalUserRequest) error
	NotAcceptExternalUser(id int) error

	// GetApplicants() ([]*dto.Applicant, error)
	// GetTerrOtdels() ([]*dto.Source, error)

	GetExternalUsers() ([]*dto.ExternalUser, error)

	// AddNewTerOtdel(terOtdel *dto.AddNewTerOtdel) error

	/// ОБНОВА

	CreateUser(user *dto.CreateUserDTO, passwordHash string) error
	GetUsers() ([]*dto.User, error)
	DeleteUser(userID string) error
	GetUserByLogin(login string) (*dto.User, error)
	GetUserProfile(userId string) (*dto.UserProfile, error)

	CreateDistrict(district *dto.District) error
	GetDistricts() ([]*dto.District, error)
	DeleteDistrict(id string) error

	CreateTerOtdel(terOtdel *dto.CreateTerOtdelDTO) (string, error)
	GetTerOtdels() ([]*dto.TerOtdel, error)
	GetTerrOtdelsByDistrictID(district_id string) ([]*dto.TerOtdel, error)
	DeleteTerOtdel(id string) error

	CreateApplicant(applicant *dto.CreateApplicantDTO) error
	GetApplicantByDistrictID(districtID string) ([]*dto.Applicant, error)

	CreateRequest(request *dto.CreateRequestDTO) error
	GetRequestsByTerOtdel(id string) ([]*dto.GetRequestsDTO, error)
	GetRequestsByDistrictID(id string) ([]*dto.GetRequestsDTO, error)
	GetRequestsByDistrictIDs(district_id string, ids []string) ([]*dto.GetRequestsDTO, error)
	ChangeStatusRequest(req *dto.ChangeStatusRequestDTO) error
	DeleteRequest(id string) error
	AddActFile(reqID string, actFile string) error

	GetRoles() ([]*dto.Role, error)
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

// func (s *Service) CreateAccount(account *dto.Account) (string, error) {
// 	err := validate.Validate(account)
// 	if err != nil {

// 		return "", ErrInvalidData
// 	}

// 	hashPass, err := security.Encode(account.Password)
// 	if err != nil {
// 		return "", err
// 	}
// 	account.Password = hashPass

// 	return s.repository.CreateAccount(account)
// }

// func (s *Service) SendRequest(request *dto.RequestFull) error {
// 	err := validate.Validate(request)
// 	if err != nil {
// 		return ErrInvalidData
// 	}
// 	fmt.Println(request, "ПРОВЕРКАА ККАСТОМАВ")

// 	fmt.Println(request.Applicant.Name, "ПРОВЕРКАА ")
// 	id, err := s.repository.NewNotPemamentApplicant(request.Applicant.Name)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("ID кастома", id)
// 	request.Applicant = dto.Applicant{
// 		ID: id,
// 	}
// 	number, err := s.repository.SendRequest(request)
// 	if err != nil {
// 		return err
// 	}

// 	if number == 0 {
// 		return ErrInvalidData
// 	}
// 	reqData := &dto.RequestForGenerating{
// 		Number:        strconv.Itoa(number),
// 		Applicant:     request.Applicant.Name,
// 		Source:        request.Source.Name,
// 		Address:       request.Address,
// 		DogsCount:     request.DogsCount,
// 		Behavior:      request.Behavior,
// 		Urgency:       request.Urgency,
// 		ContactPerson: request.ContactPerson,
// 	}

// 	if err = worker.SendRequestForGenerating(reqData); err != nil {
// 		return err
// 	}
// 	sharedDir := "/app/shared/"
// 	key := "zayavka_" + strconv.Itoa(number) + "_" + time.Now().Format("2006") + ".xlsx"
// 	filename := sharedDir + key

// 	if err = s.storage.UploadFile(context.TODO(), key, filename); err != nil {
// 		//лог
// 		fmt.Println("Error upload file to S3:", err)
// 		return err
// 	}

// 	os.Remove(filename)

// 	return nil
// }

// func (s *Service) ChangeStatusRequest(req *dto.ChangeStatusRequest) error {

// 	err := s.repository.ChangeStatusRequest(req)
// 	return err
// }

// func (s *Service) GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error) {
// 	requests, err := s.repository.GetRequestsByOtdel(otdel_id, year)
// 	if err != nil {
// 		return nil, err
// 	}
// 	go s.validateDelay(requests)
// 	return requests, nil
// }

// func (s *Service) GetAllRequests(year string) ([]*dto.RequestFull, error) {
// 	requests, err := s.repository.GetAllRequests(year)
// 	if err != nil {
// 		return nil, err
// 	}
// 	go s.validateDelay(requests)
// 	return requests, nil
// }
// func (s *Service) DeleteRequest(reqID string) error {

//		return s.repository.DeleteRequest(reqID)
//	}
// func (s *Service) ChangeProfileInfo(req *dto.ChangeProfileRequest) error {
// 	err := s.repository.ChangeProfileInfo(req)
// 	return err
// }

// func (s *Service) ChangePassword(req *dto.ChangePasswordRequest) error {
// 	var valid dto.AuthCredentials

// 	valid.Login = req.Login
// 	valid.Password = req.OldPassword

// 	_, hashPass, _, _, err := s.repository.ValidateAccount(&valid)
// 	if err != nil {
// 		return ErrInvalidData
// 	}
// 	if !security.Check(valid.Password, hashPass) {
// 		return ErrInvalidPassword
// 	}

// 	req.NewPassword, err = security.Encode(req.NewPassword)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.repository.ChangePassword(req)

// 	return err

// }

// func (s *Service) GetUserProfile(userId string) (*dto.UserProfile, error) {
// 	profile, err := s.repository.GetUserProfile(userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return profile, nil
// }

func (s *Service) GetFileURL(objectKey string) string {
	return s.storage.GetFileURL(objectKey)
}

// func (s *Service) UploadAct(req *dto.UploadActRequests, key string, filename string) error {
// 	if err := s.storage.UploadFile(context.TODO(), key, filename); err != nil {
// 		//лог
// 		fmt.Println("Error upload file to S3:", err)
// 		return err
// 	}
// 	// Удаляем временный файл
// 	os.Remove(filename)
// 	var statusReq dto.ChangeStatusRequest
// 	statusReq.ID = req.ID
// 	statusReq.Status = req.Status
// 	if err := s.ChangeStatusRequest(&statusReq); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *Service) validateDelay(requests []*dto.RequestFull) error {
// 	twoWeeks := 14 * 24 * time.Hour

// 	for _, req := range requests {
// 		if req == nil {
// 			continue
// 		}
// 		if req.Status == "Завершена" || req.Status == "Просрочена" {
// 			continue
// 		}
// 		// если прошло больше двух недель с момента создания — меняем статус
// 		if time.Since(req.CreatedAt) > twoWeeks {

// 			statusReq := &dto.ChangeStatusRequest{
// 				ID:     req.ID,
// 				Status: "Просрочена", // <-- поменяйте на нужный вам статус
// 			}

// 			if err := s.repository.ChangeStatusRequest(statusReq); err != nil {
// 				fmt.Println("failed to change status for request", req.ID, ":", err)
// 				// продолжаем обработку остальных заявок
// 			}
// 		}
// 	}

// 	return nil

// }

// func (s *Service) GenerateMultipleByDate(req *dto.GenerateMultipleRequestByDate) (string, error) {
// 	requests, err := s.repository.GetRequestsByDate(req.DateFrom, req.DateTo)
// 	if err != nil {

// 		fmt.Println("Error getting requests by number:", err)
// 		return "", err
// 	}

// 	reqs := dto.RequestForGeneratingSomething{
// 		Requests: requests,
// 		Number:   strconv.Itoa(dto.NumberOfRequests),
// 		StartRow: 20,
// 	}
// 	fmt.Println(req)
// 	if err := worker.SendRequestForGeneratingSomething(&reqs); err != nil {
// 		return "", err
// 	}

// 	sharedDir := "/app/shared/"
// 	key := "zayavka_" + reqs.Number + "_" + time.Now().Format("02-01-2006") + ".xlsx"
// 	filename := sharedDir + key

// 	if err = s.storage.UploadFile(context.TODO(), key, filename); err != nil {
// 		//лог
// 		fmt.Println("Error upload file to S3:", err)
// 		return "", err
// 	}
// 	dto.NumberOfRequests++
// 	os.Remove(filename)

// 	url := s.storage.GetFileURL(key)

// 	return url, nil
// }

// func (s *Service) GenerateMultipleByIDs(req *dto.GenerateMultipleRequestByID) (string, error) {
// 	requests, err := s.repository.GetRequestsByIDs(req.IDs)
// 	if err != nil {

// 		fmt.Println("Error getting requests by number:", err)
// 		return "", err
// 	}

// 	reqs := dto.RequestForGeneratingSomething{
// 		Requests: requests,
// 		Number:   strconv.Itoa(dto.NumberOfRequests),
// 		StartRow: 20,
// 	}

// 	fmt.Println(req)
// 	if err := worker.GenerateMultipleDocument(&reqs); err != nil {
// 		return "", err
// 	}

// 	sharedDir := "/app/shared/"
// 	key := "zayavka_" + reqs.Number + "_" + time.Now().Format("02-01-2006") + ".xlsx"
// 	filename := sharedDir + key

// 	if err = s.storage.UploadFile(context.TODO(), key, filename); err != nil {
// 		//лог
// 		fmt.Println("Error upload file to S3:", err)
// 		return "", err
// 	}
// 	dto.NumberOfRequests++
// 	os.Remove(filename)

// 	url := s.storage.GetFileURL(key)

// 	return url, nil
// }

func (s *Service) CreateExternalUser(user *dto.ExternalUser) (int, error) {

	passwordHash, err := security.Encode(user.Password)
	if err != nil {

		return 0, err
	}

	user.Password = passwordHash
	return s.repository.CreateExternalUser(user)
}

func (s *Service) ChangeSatusExternalUser(id int, status bool) error {

	return s.repository.ChangeAcceptedStatusExternalUser(id, status)

}

// Добавить внешнего чела в систему собак
func (s *Service) ChangeAcceptedExternalUser(user *dto.AcceptExternalUserRequest) error {

	return s.repository.AcceptExternalUser(user)
}

// func (s *Service) GetApplicants() ([]*dto.Applicant, error) {
// 	return s.repository.GetApplicants()
// }

// func (s *Service) GetTerrOtdels() ([]*dto.Source, error) {
// 	return s.repository.GetTerrOtdels()
// }

func (s *Service) GetExternalUsers() ([]*dto.GetExternalUser, error) {

	users, err := s.repository.GetExternalUsers()
	if err != nil {
		return nil, err
	}
	var externalUsers []*dto.GetExternalUser
	for _, user := range users {
		var usr dto.GetExternalUser
		usr.ID = user.ID
		usr.Username = user.Username
		usr.Email = user.Email
		usr.FirstName = user.FirstName
		usr.LastName = user.LastName
		usr.Patronymic = user.Patronymic
		usr.City = user.City
		usr.DateOfBirth = user.DateOfBirth
		usr.Bio = user.Bio
		usr.IsAccepted = user.IsAccepted
		externalUsers = append(externalUsers, &usr)
	}
	return externalUsers, nil
}

// func (s *Service) AddNewTerOtdel(terOtdel *dto.AddNewTerOtdel) error {
// 	return s.repository.AddNewTerOtdel(terOtdel)
// }
