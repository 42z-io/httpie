package httpie

import (
	"net/http"
	"strconv"
	"strings"
)

func GetQueryParamDefault(r *http.Request, key string, defaultValue string) string {
	ok, value := GetQueryParam(r, key)
	if !ok {
		return defaultValue
	}
	return value
}

func GetQueryParam(r *http.Request, key string) (bool, string) {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return false, ""
	}
	return true, queryValue
}

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

func GetQueryParamListDefault(r *http.Request, key string, defaultValue []string) []string {
	ok, value := GetQueryParamList(r, key)
	if !ok {
		return defaultValue
	}
	return value
}

func GetQueryParamList(r *http.Request, key string) (bool, []string) {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return false, nil
	}
	return true, strings.Split(queryValue, ",")
}
