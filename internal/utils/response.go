package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondError 返回错误响应
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

// RespondSuccess 返回成功响应
func RespondSuccess(c *gin.Context, data gin.H) {
	response := gin.H{"success": true}
	for key, value := range data {
		response[key] = value
	}
	c.JSON(http.StatusOK, response)
}
