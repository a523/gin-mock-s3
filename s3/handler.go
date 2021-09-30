package s3

import (
	"github.com/a523/gin-mock-s3/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
