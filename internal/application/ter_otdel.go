package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Create a new TerOtdel
// @Description Create a new TerOtdel
// @Accept json
// @Produce json
// @Param terOtdel body dto.CreateTerOtdelDTO true "TerOtdel"
// @Success 201 {object}  map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/ter-otdel [post]
func (h *Handlers) CreateTerOtdel(c *gin.Context) {
	var createDTO dto.CreateTerOtdelDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateTerOtdel(&createDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Get all TerOtdels
// @Description Get all TerOtdels
// @Produce json
// @Success 200 {object} []dto.TerOtdel
// @Failure 500 {object} map[string]interface{}
// @Router /api/ter-otdel [get]
func (h *Handlers) GetTerOtdels(c *gin.Context) {
	terOtdels, err := h.service.GetTerOtdels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, terOtdels)
}

// @Summary Delete a TerOtdel
// @Description Delete a TerOtdel
// @Param id path string true "TerOtdel ID"
// @Success 204 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/ter-otdel/{id} [delete]
func (h *Handlers) DeleteTerOtdel(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTerOtdel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
