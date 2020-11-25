package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignCSR(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("sign Not implemented"),
	})
}
