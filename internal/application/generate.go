package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	"github.com/niiilov/go-dog-trapping/internal/worker"
)

// @Summary Генерация документов
// @Description Генерация документов по заявкам
// @Tags Генерация
// @Accept json
// @Produce json
// @Param request body dto.Generate true "Запрос на генерацию документов"
// @Success 200 {object} dto.Response
// @Router /api/generate [post]
func (h *Handlers) Generate(c *gin.Context) {
	var req dto.Generate
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	districtID := c.GetString("district_id")

	requests, err := h.service.GetRequestsByDistrictIDs(districtID, req.RequestsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении заявок."})
		return
	}
	var reqs []*dto.RequestForGenerating

	for _, r := range requests {
		num := strconv.Itoa(r.Number)
		reqs = append(reqs, &dto.RequestForGenerating{
			Number:        num,
			ApplicantName: r.ApplicantName,
			TerOtdelName:  r.TerOtdelName,
			Address:       r.Address,
			DogsCount:     r.DogsCount,
			Behavior:      r.Behavior,
			Urgency:       r.Urgency,
			ContactPerson: r.ContactPerson,
		})
	}

	data := dto.RequestForGeneratingSomething{
		Requests: reqs,
		Number:   strconv.Itoa(dto.NumberOfRequests),
		StartRow: 20,
	}

	worker.GenerateMultipleDocument(c, &data)

}
