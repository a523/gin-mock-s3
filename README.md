# Gin-mock-s3

mock api from:

https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html
http://docs.ceph.org.cn/radosgw/s3/objectops/
https://help.aliyun.com/document_detail/31947.html

## env
```shell
DRIVER # storage type. 'local' or 's3'. default 'local'
BASIC_PATH # When store at local, the local path. default '/tmp'
S3_ENDPOINT # When store with s3, It's need. default 'http://127.0.0.1:9000'
S3_ACCESS_KEY # When store with s3, It's need. default 'minioadmin'
S3_SECRET_KEY # When store with s3, It's need. default 'minioadmin'
```
## run
```shell
go run .
```

