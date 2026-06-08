package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	jw "github.com/niiilov/go-dog-trapping/pkg/jwt"
)

type Service interface {
	// 	CreateAccount(account *dto.Account) (string, error)

	// 	SendRequest(request *dto.RequestFull) error
	// 	DeleteRequest(reqID string) error
	// 	UploadAct(req *dto.UploadActRequests, key string, filename string) error

	// 	ChangeStatusRequest(req *dto.ChangeStatusRequest) error
	// //	ChangePassword(req *dto.ChangePasswordRequest) error
	// //	ChangeProfileInfo(req *dto.ChangeProfileRequest) error

	// 	GenerateMultipleByDate(req *dto.GenerateMultipleRequestByDate) (string, error)
	// 	GenerateMultipleByIDs(req *dto.GenerateMultipleRequestByID) (string, error)

	// GetAllRequests(year string) ([]*dto.RequestFull, error)
	// GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error)
	// GetUserProfile(userId string) (*dto.UserProfile, error)
	GetFileURL(objectKey string) string

	CreateExternalUser(user *dto.ExternalUser) (int, error)
	ChangeSatusExternalUser(id int, status bool) error

	ChangeAcceptedExternalUser(user *dto.AcceptExternalUserRequest) error

	// GetApplicants() ([]*dto.Applicant, error)
	// GetTerrOtdels() ([]*dto.Source, error)

	// GetExternalUsers() ([]*dto.GetExternalUser, error)

	// AddNewTerOtdel(terOtdel *dto.AddNewTerOtdel) error

	// Обнова

	CreateUser(user *dto.CreateUserDTO) error
	ValidateAccount(account *dto.AuthCredentials) (profile *dto.UserProfile, err error)
	GetUsers() ([]*dto.GetUserDTO, error)
	DeleteUser(userID string) error

	CreateTerOtdel(terOtdel *dto.CreateTerOtdelDTO) (string, error)
	GetTerOtdels() ([]*dto.TerOtdel, error)
	DeleteTerOtdel(id string) error

	GetApplicantByTerOtdelID(terOtdelID string) ([]*dto.Applicant, error)
	CreateApplicant(applicant *dto.CreateApplicantDTO) error

	CreateRequest(request *dto.CreateRequestDTO) error
	GetAllRequests() ([]*dto.GetRequestsDTO, error)
	GetRequestsByTerOtdel(id string) ([]*dto.GetRequestsDTO, error)
	GetRequestsByIDs(ids []string) ([]*dto.GetRequestsDTO, error)
	GetRequestsByTerOtdelIDs(terOtdelID string, ids []string) ([]*dto.GetRequestsDTO, error)
	ChangeStatusRequest(req *dto.ChangeStatusRequestDTO) error
	DeleteRequest(id string) error
	AddActFile(reqID string, actFilename, actFilePath string) error

	GetRoles() ([]*dto.Role, error)
}
type Handlers struct {
	jwtService *jw.ServiceJWT

	service Service
}

func NewHandlers(service Service, jwtService *jw.ServiceJWT) *Handlers {
	return &Handlers{service: service, jwtService: jwtService}
}

// @Summary Download Act File URL
// @Security BearerAuth
// @Description Получение URL для скачивания файла акта
// @Tags Requests
// @Produce json
// @Param filename path string true "Filename"
// @Success 200 {object} dto.ResponseUrl
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении запросов."
// @Router /api/requests/act/{filename} [get]
func (h *Handlers) DownloadAct(c *gin.Context) {

	filename := c.Param("filename")
	url := h.service.GetFileURL(filename)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url})
}

// Создать внешнего пользователя POST
// @Summary Create external user
// @Security BearerAuth
// @Description Create external user
// @Tags external-users
// @Accept json
// @Produce json
// @Param user body dto.ExternalUser true "External user to create"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/external/users [post]
func (h *Handlers) CreateExternalUser(c *gin.Context) {
	var user dto.ExternalUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateExternalUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary ChangeAccept external user
// @Security BearerAuth
// @Description Change Accepted external user
// @Tags external-users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/external/users/{id}/activate [put]
func (h *Handlers) ChangeStatusAcceptExternalUser(c *gin.Context) {
	strID := c.Param("id")

	id, err := strconv.Atoi(strID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeSatusExternalUser(id, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "External user status updated successfully"})
}

// @Summary ChangeDisAccept external user
// @Security BearerAuth
// @Description Change Accepted external user
// @Tags external-users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/external/users/{id} [put]
func (h *Handlers) ChangeStatusDisAcceptExternalUser(c *gin.Context) {
	strID := c.Param("id")

	id, err := strconv.Atoi(strID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeSatusExternalUser(id, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "External user status updated successfully"})
}

// @Summary Accept internal systen external user
// @Security BearerAuth
// @Description  Accepted external user to dogs
// @Tags external-users
// @Accept json
// @Produce json
// @Param user body dto.AcceptExternalUserRequest true "External user to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/external/{id} [put]
func (h *Handlers) ChangeAcceptedExternalUser(c *gin.Context) {
	strID := c.Param("id")

	id, err := strconv.Atoi(strID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req dto.AcceptExternalUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.service.ChangeAcceptedExternalUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "External user status updated successfully"})
}
