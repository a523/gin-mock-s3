package main

import (
	"github.com/a523/gin-mock-s3/config"
	"strconv"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080

	addr := ":" + strconv.Itoa(config.CFG.Port)
	err := r.Run(addr)
	panic(err)
}
