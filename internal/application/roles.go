package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all roles
// @Security BearerAuth
// @Description Получение всех ролей
// @Tags roles
// @Produce json
// @Success 200 {array} dto.Role
// @Failure 500 {object} dto.Response	"Ошибка при получении ролей."
// @Router /api/roles [get]
func (h *Handlers) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}
