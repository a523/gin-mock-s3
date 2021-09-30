package main

import (
	"github.com/a523/gin-mock-s3/config"
	"github.com/a523/gin-mock-s3/s3"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	if config.CFG.Driver == "local" {
		router.GET("/:bucket/", getBucket)
		router.HEAD("/:bucket/", headBucket)
		router.POST("/:bucket/", deleteObjects)
		router.PUT("/:bucket/:object/*objectSuffix", uploadFile)
		router.DELETE("/:bucket/:object/*objectSuffix", deleteFile)
		router.GET("/:bucket/:object/*objectSuffix", getFile)
	} else if config.CFG.Driver == "s3" {
		router.GET("/:bucket/", s3.ListObjects)
		router.HEAD("/:bucket/", s3.HeadBucket)
		router.POST("/:bucket/", s3.DeleteObjects)
		router.PUT("/:bucket/:object/*objectSuffix", s3.UploadObject)
		router.GET("/:bucket/:object/*objectSuffix", s3.GetObject)
	} else {
		panic("Not supported the driver " + config.CFG.Driver)
	}

	return router
}
