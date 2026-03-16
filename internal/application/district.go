package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Get all districts
// @Description Get all districts
// @Tags districts
// @Produce json
// @Success 200 {object} []dto.District
// @Failure 500 {object} map[string]string
// @Router /api/districts [get]
func (h *Handlers) GetDistricts(c *gin.Context) {
	districts, err := h.service.GetDistricts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, districts)
}

// @Summary Create a district
// @Description Create a new district
// @Tags districts
// @Accept json
// @Produce json
// @Param district body dto.District true "District"
// @Success 201 {object} dto.District
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/districts [post]
func (h *Handlers) CreateDistrict(c *gin.Context) {
	var district dto.District
	if err := c.ShouldBindJSON(&district); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateDistrict(&district); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Успешное создание МО")
}

// @Summary Delete a district
// @Description Delete a district by ID
// @Tags districts
// @Produce json
// @Param id path string true "District ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/districts/{id} [delete]
func (h *Handlers) DeleteDistrict(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteDistrict(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
