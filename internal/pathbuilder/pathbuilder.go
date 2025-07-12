package pathbuilder

import (
	"errors"
	"regexp"
	"strings"
)

func BuildTargetURL(baseURL, path string) (string, error) {

	path = strings.Trim(strings.TrimPrefix(path, "/api/"), "/")
	if path == "" {
		return "", errors.New("empty path after /api/")
	}

	parts := strings.Split(path, "/")
	if parts[0] != "posts" {
		return "", errors.New("invalid endpoint")
	}

	targetURL := baseURL + "/posts"

	if len(parts) == 2 {
		if matched, _ := regexp.MatchString(`^\d+$`, parts[1]); !matched {
			return "", errors.New("invalid post ID")
		}
		targetURL += "/" + parts[1]
	} else if len(parts) > 2 {
		return "", errors.New("invalid path")
	}

	return targetURL, nil
}
