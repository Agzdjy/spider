package util

import (
	"os"
	"fmt"
	"strings"
)

func Exist(fileName string) bool {
	fi, err := os.Stat(fileName)

	fmt.Println(fi, err)

	return err == nil || os.IsExist(err)
}

func GetEnv(key string) string {
	return strings.Trim(os.Getenv(key), " ")
}