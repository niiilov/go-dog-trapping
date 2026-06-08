package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Get Applicants by TerOtdel ID
// @Description Получение заявителей по ID территориального отдела
// @Tags Applicants
// @Accept json
// @Produce json
// @Param terOtdelID path string true "TerOtdel ID"
// @Success 200 {array} dto.Applicant	"Заявители по территориальному отделу"
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении заявителей."
// @Router /api/applicants/{terOtdelID} [get]
func (h *Handlers) GetApplicantByTerOtdelID(c *gin.Context) {
	terOtdelID := c.Param("terOtdelID")

	applicants, err := h.service.GetApplicantByTerOtdelID(terOtdelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении заявителей."})
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
