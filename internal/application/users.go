package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Получить список пользователей
// @Description Получить список всех пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetUserDTOResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/users [get]
func (h *Handlers) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении пользователей."})
		return
	}
	c.JSON(http.StatusOK, dto.GetUserDTOResponse{Status: "ok", Users: users})
}

// @Summary Удалить пользователя
// @Description Удалить пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/{id} [delete]
func (h *Handlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при удалении пользователя."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Пользователь успешно удален."})
}
