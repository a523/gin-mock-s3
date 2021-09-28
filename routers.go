package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.GET("/:bucket/", getBucket)
	router.HEAD("/:bucket/", headBucket)
	router.POST("/:bucket/", deleteObjects)
	router.PUT("/:bucket/:object/*objectSuffix", uploadFile)
	router.DELETE("/:bucket/:object/*objectSuffix", deleteFile)
	router.GET("/:bucket/:object/*objectSuffix", getFile)

	return router
}
