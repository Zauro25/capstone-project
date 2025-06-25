package utils

import "github.com/gin-gonic/gin"

// SuccessResponse mengirimkan respons sukses
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// ErrorResponse mengirimkan respons error
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"error":   err,
	})
}
