package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func GenerateMultipleDocument(c *gin.Context, reqData *dto.RequestForGeneratingSomething) {
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to marshal request: %v", err)})
		return
	}
	fmt.Println("Sending request with data:", string(jsonData))

	client := &http.Client{}
	resp, err := client.Post("http://worker-sobaki:8001/generate-multiple", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("ЭТО ТЕСТ", bytes.NewReader(jsonData))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	// Пробрасываем заголовки от питона (Content-Disposition, Content-Type)
	for key, values := range resp.Header {
		for _, v := range values {
			c.Header(key, v)
		}
	}

	c.Status(resp.StatusCode)

	// Стримим тело — не грузим весь файл в память
	io.Copy(c.Writer, resp.Body)

}
