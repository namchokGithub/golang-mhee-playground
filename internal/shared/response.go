package shared

import "github.com/gin-gonic/gin"

type APIResponse struct {
	OK    bool      `json:"ok"`
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, APIResponse{OK: true, Data: data})
}

func Fail(c *gin.Context, httpStatus int, code, msg string) {
	c.JSON(httpStatus, APIResponse{
		OK: false,
		Error: &APIError{
			Code:    code,
			Message: msg,
		},
	})
}
