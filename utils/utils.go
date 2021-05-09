package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func FindInSlice(slice []string, val string) (bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetBooleanEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	booleanValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Error(fmt.Sprintf("%s should be a boolean",key))
		panic(err)
	}
	return booleanValue
}

func IsIpv4WithPort(host string) bool {
	portPart := strings.Split(host,":")
	parts := strings.Split(portPart[0], ".")

	if len(parts) < 4 {
		return false
	}

	for _,x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}

	if len(portPart) != 2 {
		return false
	}

	return true
}