package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/a523/gin-mock-s3/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func InitClient() *s3.Client {
	endpoint := config.CFG.S3Endpoint
	accessKey := config.CFG.S3AccessKey
	secretKey := config.CFG.S3SecretKey

	if endpoint == "" || accessKey == "" || secretKey == "" {
		panic("Not find minio config")
	}

	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg := aws.Config{
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		EndpointResolver: staticResolver,
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
		options.EndpointOptions.DisableHTTPS = true
		options.Region = "us-east-1"
	})

	return client
}

func ListObjects(c *gin.Context) {
	_, queryLocation := c.GetQuery("location")
	bucket := c.Param("bucket")
	s3Client := InitClient()
	if queryLocation {
		lInput := s3.GetBucketLocationInput{Bucket: &bucket}
		out, err := s3Client.GetBucketLocation(context.TODO(), &lInput)
		if err != nil {
			var re *awshttp.ResponseError
			if errors.As(err, &re) {
				c.XML(re.HTTPStatusCode(), err.Error())
			} else {
				c.XML(http.StatusInternalServerError, err.Error())
			}

		} else {
			c.XML(http.StatusOK, &out)
		}
		return
	}
	listObjectsInput := s3.ListObjectsInput{Bucket: &bucket}
	err := c.BindQuery(listObjectsInput)
	if err != nil {
		c.XML(http.StatusBadRequest, err.Error())
		return
	}

	out, err := s3Client.ListObjects(context.TODO(), &listObjectsInput)
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			c.XML(re.HTTPStatusCode(), err.Error())
			return
		}
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	c.XML(http.StatusOK, &out)
}

func HeadBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	headBucketInput := s3.HeadBucketInput{Bucket: &bucket}
	err := c.BindQuery(headBucketInput)
	if err != nil {
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	s3Client := InitClient()

	out, err := s3Client.HeadBucket(context.TODO(), &headBucketInput)
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			c.XML(re.HTTPStatusCode(), err.Error())
			return
		}
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	c.XML(http.StatusOK, &out)
}

func DeleteObjects(c *gin.Context) {
	bucket := c.Param("bucket")
	_, del := c.GetQuery("delete")

	if del {
		deleteObjectsInput := s3.DeleteObjectsInput{Bucket: &bucket}
		err := c.BindXML(deleteObjectsInput)
		if err != nil {
			c.XML(http.StatusBadRequest, err.Error())
			return
		}
		err = c.BindQuery(deleteObjectsInput)

		s3Client := InitClient()

		out, err := s3Client.DeleteObjects(context.TODO(), &deleteObjectsInput)
		if err != nil {
			var re *awshttp.ResponseError
			if errors.As(err, &re) {
				c.XML(re.HTTPStatusCode(), err.Error())
				return
			}
			c.XML(http.StatusBadRequest, err.Error())
			return
		}
		c.XML(http.StatusOK, &out)
		return
	} else {
		c.XML(http.StatusInternalServerError, "not supported")
		return
	}
}

func UploadObject(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectSuffix := c.Param("objectSuffix")
	key := path.Join(object, objectSuffix)

	data, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error uploading : %s", err)
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	putObjectInput := s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: bytes.NewReader(data)}
	err = c.BindQuery(putObjectInput)
	if err != nil {
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	s3Client := InitClient()

	out, err := s3Client.PutObject(context.TODO(), &putObjectInput)
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			c.XML(re.HTTPStatusCode(), err.Error())
			return
		}
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	c.XML(http.StatusOK, &out)
}

func GetObject(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	objectSuffix := c.Param("objectSuffix")
	key := path.Join(object, objectSuffix)
	getObjectInput := s3.GetObjectInput{Bucket: &bucket, Key: &key}
	s3Client := InitClient()

	out, err := s3Client.GetObject(context.TODO(), &getObjectInput)
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			c.XML(re.HTTPStatusCode(), err.Error())
			return
		}
		c.XML(http.StatusBadRequest, err.Error())
		return
	}
	extraHeaders := map[string]string{
		"Last-Modified": out.LastModified.UTC().Format(http.TimeFormat),
		"ETag":          *out.ETag,
	}
	c.DataFromReader(200, out.ContentLength, *out.ContentType, out.Body, extraHeaders)
}
