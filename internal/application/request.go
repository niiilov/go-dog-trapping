package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Get Requests
// @Security BearerAuth
// @Description Получение заявок на отлов бродячих собак по ID территориального отдела.
// @Tags Requests
// @Produce json
// @Success 200 {object} []dto.Request	"Список заявок"
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении заявок."
// @Router /api/requests [get]
func (h *Handlers) GetRequests(c *gin.Context) {
	roleID := c.GetString("role")
	terOtdelID := c.GetString("ter_otdel_id")

	var requests []*dto.GetRequestsDTO
	var err error

	if roleID == dto.RoleRegionalAdmin {
		requests, err = h.service.GetAllRequests()
	} else {
		requests, err = h.service.GetRequestsByTerOtdel(terOtdelID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type nestedRef struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type requestResponse struct {
		ID            string    `json:"id"`
		Number        string    `json:"number"`
		Address       string    `json:"address"`
		DogsCount     int       `json:"dogs_count"`
		Behavior      string    `json:"behavior"`
		Urgency       string    `json:"urgency"`
		ContactPerson string    `json:"contact_person"`
		Status        string    `json:"status"`
		CreatedAt     string    `json:"created_at"`
		ActFile       string    `json:"act_file,omitempty"`
		Applicant     nestedRef `json:"applicant"`
		Source        nestedRef `json:"source"`
	}

	var response []requestResponse
	for _, r := range requests {
		response = append(response, requestResponse{
			ID:            r.ID,
			Number:        fmt.Sprintf("%d", r.Number),
			Address:       r.Address,
			DogsCount:     r.DogsCount,
			Behavior:      r.Behavior,
			Urgency:       r.Urgency,
			ContactPerson: r.ContactPerson,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			ActFile:       r.ActFile,
			Applicant:     nestedRef{ID: r.ApplicantID, Name: r.ApplicantName},
			Source:        nestedRef{ID: r.TerOtdelID, Name: r.TerOtdelName},
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Create Request
// @Description Создание заявки на отлов бродячих собак.
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body dto.CreateRequestDTO true "Создание заявки"
// @Success 201 {object} dto.Response "Заявка успешно создана"
// @Failure 400 {object} dto.Response "Ошибка в данных запроса."
// @Failure 500 {object} dto.Response "Ошибка при создании заявки."
// @Router /api/requests [post]
func (h *Handlers) CreateRequest(c *gin.Context) {
	var req dto.CreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRequest(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully"})
}

// @Summary Change Request Status
// @Security BearerAuth
// @Description Изменение статуса заявки на отлов бродячих собак.
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body dto.ChangeStatusRequestDTO true "Изменение статуса заявки"
// @Success 200 {object} dto.Response "Статус заявки успешно изменен"
// @Failure 400 {object} dto.Response "Ошибка в данных запроса."
// @Failure 500 {object} dto.Response "Ошибка при изменении статуса заявки."
// @Router /api/requests/status [put]
func (h *Handlers) ChangeStatusRequest(c *gin.Context) {
	var req dto.ChangeStatusRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeStatusRequest(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully"})
}

// @Summary Delete Request
// @Security BearerAuth
// @Description Удаление заявки на отлов бродячих собак.
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} dto.Response "Заявка успешно удалена"
// @Failure 400 {object} dto.Response "Ошибка в данных запроса."
// @Failure 500 {object} dto.Response "Ошибка при удалении заявки."
// @Router /api/requests/{id} [delete]
func (h *Handlers) DeleteRequest(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}

// @Summary Add Act File
// @Security BearerAuth
// @Description Добавление акта об отлове бродячих собак.
// @Tags Requests
// @Accept json
// @Produce json
// @Param  file formData file true "Act File"
// @Param request_id formData string true "Request ID"
// @Success 200 {object} dto.Response "Акт успешно добавлен"
// @Failure 400 {object} dto.Response "Ошибка в данных запроса."
// @Failure 500 {object} dto.Response "Ошибка при добавлении акта."
// @Router /api/requests/act [post]
func (h *Handlers) AddActFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	filePath := "/app/act/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}
	reqID := c.PostForm("request_id")
	if err := h.service.AddActFile(reqID, file.Filename, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Act file added successfully"})
}
