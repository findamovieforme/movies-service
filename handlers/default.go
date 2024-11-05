// handlers/movie.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) Default(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Movies service started successfully!",
	})
}

func (h *DefaultHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
