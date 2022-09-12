package main

import (
	"errors"
	"net/http"
	"strings"
)

var errNoAuthHeaderFound = errors.New("empty or no authorization header found")
var errAuthHeaderWrongFormat = errors.New("authorization header in wrong format")
var errEmptyToken = errors.New("empty bearer token found")

func getAuthToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errNoAuthHeaderFound
	}

	splits := strings.SplitN(token, " ", 2)
	if len(splits) < 2 {
		return "", errAuthHeaderWrongFormat
	}

	if splits[0] != "Bearer" {
		return "", errAuthHeaderWrongFormat
	}

	if splits[1] == "" {
		return "", errEmptyToken
	}

	return splits[1], nil
}
