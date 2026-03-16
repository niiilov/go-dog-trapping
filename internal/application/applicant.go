package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Get Applicants by District ID
// @Description Получение заявок по ID района
// @Tags Applicants
// @Accept json
// @Produce json
// @Param districtID path string true "District ID"
// @Success 200 {array} dto.Applicant	"Заявки по району"
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении заявок."
// @Router /api/applicants/{districtID} [get]
func (h *Handlers) GetApplicantByDistrictID(c *gin.Context) {
	districtID := c.Param("districtID")

	applicants, err := h.service.GetApplicantByDistrictID(districtID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении заявок."})
		return
	}

	c.JSON(http.StatusOK, applicants)
}

// @Summary Create Applicant
// @Description Создание заявки на отлов бродячей собаки
// @Tags Applicants
// @Accept json
// @Produce json
// @Param request body dto.CreateApplicantDTO true "Applicant information"
// @Success 200 {object} dto.Response	"Заявка успешно создана."
// @Failure 400 {object} dto.Response  	"Ошибка в данных заявки."
// @Failure 500 {object} dto.Response	"Ошибка при создании заявки."
// @Router /api/applicants [post]
func (h *Handlers) CreateApplicant(c *gin.Context) {
	var applicant dto.CreateApplicantDTO
	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных заявки."})
		return
	}

	if err := h.service.CreateApplicant(&applicant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при создании заявки."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Заявитель успешно создана."})
}
