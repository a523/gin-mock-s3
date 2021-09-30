package main

import (
	"fmt"
	"github.com/a523/gin-mock-s3/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

func uploadFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectSuffix := c.Param("objectSuffix")
	objectPath := path.Join(config.CFG.BasicPath, bucket, object, objectSuffix)
	objectDir := path.Dir(objectPath)
	err := os.MkdirAll(objectDir, os.ModePerm)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	data, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error uploading : %s", err)
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	err = ioutil.WriteFile(objectPath, data, 0644)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	md5 := getFileMd5(objectPath)
	c.XML(http.StatusOK, putObjectResponse{Etag: md5})
	return
}

func deleteFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	//objectSuffix := c.Param("objectSuffix")
	objectBasicDir := path.Join(config.CFG.BasicPath, bucket, object)
	err := os.RemoveAll(objectBasicDir)
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	c.XML(http.StatusOK, nil)
	return
}

func deleteObjects(c *gin.Context) {
	bucket := c.Param("bucket")
	_, del := c.GetQuery("delete")
	if del {
		bucketPath := path.Join(config.CFG.BasicPath, bucket)
		dir, _ := ioutil.ReadDir(bucketPath)
		for _, d := range dir {
			err := os.RemoveAll(path.Join([]string{bucketPath, d.Name()}...))
			if err != nil {
				fmt.Printf("Failed to remove objects:%s", err)
				c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
				return
			}
		}

		c.XML(http.StatusOK, nil)
		return
	}
}

func getFile(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectSuffix := c.Param("objectSuffix")
	objectPath := path.Join(config.CFG.BasicPath, bucket, object, objectSuffix)
	//fileInfo, err := os.Stat(objectPath)

	data, err := os.Open(objectPath)
	fileInfo, err := data.Stat()
	if err != nil {
		c.XML(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	extraHeaders := map[string]string{
		"Last-Modified": fileInfo.ModTime().UTC().Format(http.TimeFormat),
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`,
			url.QueryEscape(fileInfo.Name()), url.QueryEscape(fileInfo.Name()),
		),
		//"ETag": *resp.ETag,
	}

	c.DataFromReader(http.StatusOK, fileInfo.Size(), "application/octet‑stream", data, extraHeaders)
}

func getBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	_, exists := c.GetQuery("location")
	if exists {
		c.XML(http.StatusOK, "")
		return
	}

	bucketPath := path.Join(config.CFG.BasicPath, bucket)
	files, _ := ioutil.ReadDir(bucketPath)

	var contents []Content
	for _, f := range files {
		contents = append(contents, Content{Size: f.Size(), Key: f.Name(), LastModified: f.ModTime()})
	}
	result := ListBucketResult{Name: bucket, Contents: contents}
	c.XML(http.StatusOK, result)
}

func headBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	bucketPath := path.Join(config.CFG.BasicPath, bucket)
	if isDir(bucketPath) {
		c.XML(http.StatusOK, "")
	} else {
		fmt.Print("can't find bucket")
		c.XML(http.StatusConflict, "")
	}
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
