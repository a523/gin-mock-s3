package main

import (
	"os"
)

var CFG Config

type Config struct {
	BasicPath   string
	Driver      string // local or s3
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
}

var defConf = Config{
	BasicPath:   "/tmp/",
	Driver:      "local",
	S3Endpoint:  "http://127.0.0.1:9000",
	S3AccessKey: "minioadmin",
	S3SecretKey: "minioadmin",
}

func getConfig() Config {
	var cfg Config
	cfg.BasicPath = readEnvConf("BASIC_PATH", defConf.BasicPath)
	cfg.Driver = readEnvConf("DRIVER", defConf.Driver)
	cfg.S3Endpoint = readEnvConf("S3_ENDPOINT", defConf.S3Endpoint)
	cfg.S3AccessKey = readEnvConf("S3_ACCESS_KEY", defConf.S3AccessKey)
	cfg.S3SecretKey = readEnvConf("S3_SECRET_KEY", defConf.S3SecretKey)
	return cfg
}

func readEnvConf(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	} else {
		return def
	}
}

func checkConf(cfg Config) {
	if cfg.Driver != "local" && cfg.Driver != "s3" {
		panic(`The config 'DRIVER' must be 'local' or 's3'`)
	}
}

func init() {
	CFG = getConfig()
	checkConf(CFG)
}
