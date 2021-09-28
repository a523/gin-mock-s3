package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func uploadFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectPath := path.Join(BasicPath, bucket, object)
	data, err := c.GetRawData()
	if err != nil {
		c.XML(http.StatusBadRequest, nil)
		return
	}
	err = ioutil.WriteFile(objectPath, data, 0644)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	c.XML(http.StatusOK, nil)
	return
}

func deleteFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectPath := path.Join(BasicPath, bucket, object)
	err := os.Remove(objectPath)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	c.XML(http.StatusOK, nil)
	return
}

func getFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectPath := path.Join(BasicPath, bucket, object)
	data, err := ioutil.ReadFile(objectPath)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	c.Stream(func(w io.Writer) bool {
		n, err := w.Write(data)
		if err != nil {
			return false
		} else {
			c.Header("ContentLength", string(rune(n)))
			return true
		}

	})
}
