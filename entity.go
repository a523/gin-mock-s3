package main

type Response struct {
	ContentLength string `xml:"Content-Length"`
	Etag          string `xml:"Etag"`
}

type Request struct {
	ContentMd5  string `xml:"content-md5"`
	ContentType string `xml:"content-type"`
}

type ObjectResponse struct {
	Size int64 `xml:"Size"`
}

type ErrorResponse struct {
	Err string `xml:"Error"`
}
