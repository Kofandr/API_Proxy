package utils

import (
	"net/http"
	"strings"
)

func GetSupportedMethods(targetURL string) []string {
	if strings.Contains(targetURL, "/posts/") {
		return []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
	}
	return []string{http.MethodGet, http.MethodPost}
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
