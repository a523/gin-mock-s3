package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func getFileMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("failed to open file %s: %v", filePath, err)
		return ""
	}

	defer file.Close()
	md5 := md5.New()
	io.Copy(md5, file)
	return hex.EncodeToString(md5.Sum(nil))
}
