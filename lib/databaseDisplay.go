package lib

import (
	"strings"
)

func GetDisplayURL(url string) string {
	parts := strings.Split(url, ":")
	password := strings.Split(parts[2], "@")[0]
	return strings.Replace(url, password, "********", 1)
}
