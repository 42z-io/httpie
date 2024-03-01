package httpie

import (
	"net/http"
	"strconv"
	"strings"
)

// Get a query parameter (string) from the request, use a default value if it does not exist
func GetQueryParamDefault(r *http.Request, key string, defaultValue string) string {
	ok, value := GetQueryParam(r, key)
	if !ok {
		return defaultValue
	}
	return value
}

// Get a query parameter (string) from the request, return a boolean if it exists
func GetQueryParam(r *http.Request, key string) (bool, string) {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return false, ""
	}
	return true, queryValue
}

// Get a query parameter (int) from the request, use a default value if it does not exist and an error if its invalid
func GetQueryParamIntDefault(r *http.Request, key string, defaultValue int) (int, error) {
	ok, value, err := GetQueryParamInt(r, key)
	if err != nil {
		return 0, err
	}
	if !ok {
		return defaultValue, nil
	}
	return value, nil
}

// Get a query parameter (int) from the request, return a boolean if it exists, and an error if its invalid
func GetQueryParamInt(r *http.Request, key string) (bool, int, error) {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return false, 0, nil
	}
	value, err := strconv.Atoi(queryValue)
	if err != nil {
		return true, 0, err
	}
	return true, value, nil
}

// Get a query parameter ([]string) from the request, split the value by commas, use a default value if it does not exist
func GetQueryParamListDefault(r *http.Request, key string, defaultValue []string) []string {
	ok, value := GetQueryParamList(r, key)
	if !ok {
		return defaultValue
	}
	return value
}

// Get a query parameter ([]string) from the request, split the value by commas, return a boolean if it exists
func GetQueryParamList(r *http.Request, key string) (bool, []string) {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return false, nil
	}
	return true, strings.Split(queryValue, ",")
}
