package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadFile(c *gin.Context) {
	c.XML(http.StatusOK, nil)
	return
}
