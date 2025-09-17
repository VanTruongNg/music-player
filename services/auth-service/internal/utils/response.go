package utils

import "github.com/gin-gonic/gin"

// Response is the standard API response wrapper.
type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Meta  interface{}    `json:"meta,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse is the standard API error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success sends a standard success response.
func Success(c *gin.Context, status int, data interface{}, meta ...interface{}) {
	resp := Response{Data: data}
	if len(meta) > 0 {
		resp.Meta = meta[0]
	}
	c.JSON(status, resp)
}

// Fail sends a standard error response.
func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}
