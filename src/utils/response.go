package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is a standard JSON response structure.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success sends a standard success response with data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMessage sends a standard success response with a custom message and no data.
func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  message,
		Data: nil,
	})
}

// Fail sends a generic error response.
func Fail(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{ // Often business APIs return 200 OK with an error code inside
		Code: -1,
		Msg:  message,
		Data: nil,
	})
}

// FailWithBadRequest sends a 400 Bad Request response.
func FailWithBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: -1,
		Msg:  message,
		Data: nil,
	})
}

// FailWithUnauthorized sends a 401 Unauthorized response.
func FailWithUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: -1,
		Msg:  message,
		Data: nil,
	})
}
