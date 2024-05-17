package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
func badRequest(c *gin.Context, message ...string) {
	if len(message) == 0 {
		message = []string{"Bad Request"}
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": message[0]})
}
