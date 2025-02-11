package utils

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// JSON success response model
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BuildResponse(c *fiber.Ctx, status int, message string, data interface{}) {
	result := SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.Status(status).JSON(result)
}

func BuildErrorResponse(c *fiber.Ctx, status int, message string) {
	c.Response().Header.Set("Content-Type", "application/json")
	c.Status(status).JSON(ErrorResponse{Status: status, Message: message})
}

// func BuildErrorResponse(c *gin.Context, status int, message string) {
// 	c.Writer.Header().Set("Content-Type", "application/json")
// 	c.IndentedJSON(status, ErrorResponse{Status_code: status, Status_message: message})
// }

// func BuildErrorResponse(r http.ResponseWriter, status int, message string) {
// 	result := ErrorResponse{
// 		Status_code:    status,
// 		Status_message: message,
// 	}
// 	r.Header().Set("Content-Type", "application/json")
// 	r.WriteHeader(status)
// 	json.NewEncoder(r).Encode(result)
// }
