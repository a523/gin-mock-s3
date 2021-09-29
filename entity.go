package main

import "time"

type ErrorResponse struct {
	Err string `xml:"Error"`
}

type ListBucketResult struct {
	Name           string    `xml:"Name"`
	Prefix         string    `xml:"Prefix"`
	Marker         string    `xml:"Marker"`
	MaxKeys        int       `xml:"MaxKeys"`
	Delimiter      string    `xml:"Delimiter"`
	IsTruncated    bool      `xml:"IsTruncated"`
	CommonPrefixes []string  `xml:"CommonPrefixes"`
	Contents       []Content `xml:"Contents"`
}

type Content struct {
	Key          string    `xml:"key"`
	LastModified time.Time `xml:"LastModified"`
	Etag         string    `xml:"Etag"`
	Size         int64     `xml:"Size"`
}

type putObjectResponse struct {
	Etag         string    `xml:"Etag"`
}