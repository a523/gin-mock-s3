package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.PUT("/:bucket/:object", uploadFile)
	router.DELETE("/:bucket/:object", deleteFile)
	router.GET("/:bucket/:object", getFile)
	return router
}
