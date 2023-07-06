package Utility

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

var Variables []string = []string{
	"GOCHAT_USER",
	"GOCHAT_PASSWORD",
	"GOCHAT_DBNAME",
	"GOCHAT_SSLMODE",
	"GOCHAT_SOCKET",
}

func GetConfig() (error, map[string]string) {
	config := make(map[string]string)
	for _, key := range Variables {
		value, exist := os.LookupEnv(key)
		if !exist {
			return errors.New(fmt.Sprintf("Environment variable %s isn't setted correctly", key)), config
		}
		config[strings.ToLower(strings.Split(key, "_")[1])] = value
	}

	return nil, config
}
