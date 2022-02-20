package handler

import (
	"net/url"
	"regexp"
	"strconv"
)

// getFirstPathValue searches for the first submatch in of the provided regex pattern
// inside the given path string and returns it. returns an empty string if no match was found
func getFirstPathValue(pattern string, path string) string {
	regex := regexp.MustCompile(pattern)
	regexResult := regex.FindStringSubmatch(path)

	if regexResult == nil {
		return ""
	}

	if len(regexResult) < 2 {
		return ""
	}

	return regexResult[1]
}

func getStringParam(queryParams url.Values, param string, defaultValue string) string {
	if val, ok := queryParams[param]; ok && len(val) > 0 {
		return val[0]
	}
	return defaultValue
}

func getIntParam(queryParams url.Values, param string, defaultValue int) int {
	if val, ok := queryParams[param]; ok && len(val) > 0 {
		i, err := strconv.Atoi(val[0])
		if err != nil  {
			return defaultValue
		}

		return i
	}
	return defaultValue
}