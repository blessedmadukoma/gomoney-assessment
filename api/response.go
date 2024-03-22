package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseMessage struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// successResponse handles the success response by using map[string]interface{} to return the message and data
// func successResponse(s string, data interface{}) gin.H {
func successResponse(s string, data interface{}) ResponseMessage {
	response := ResponseMessage{
		Message: s,
		Data:    data,
	}
	// return gin.H{"message": s, "data": data}
	return response
}

// errorResponse handles the error response by using map[string]interface{} to return the error and it's message
// func errorResponse(s string, err error) gin.H {
func errorResponse(s string, err error) ResponseMessage {
	if err != nil {
		response := ResponseMessage{
			Message: fmt.Sprintf(s + " -> " + err.Error()),
			Data:    nil,
		}

		fmt.Println("response:", response)

		// return gin.H{"error: ": s + " -> " + err.Error()}
		return response
	}

	response := ResponseMessage{
		Message: s,
		Data:    nil,
	}
	// return gin.H{"error": s}
	return response
}

// rateLimitExceededResponse is a helper to send  a 429 To Many Requests
func (srv *Server) rateLimitExceededResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, errorResponse("rate limit exceeded", nil))
}

type ErrDocumentation struct {
	ErrorMessage string
}
