package utils

import (
	"net/url"
)

func AddPrefix(addr string) (string, error) {

	v, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	if v.Scheme == "" {
		v.Scheme = "http"
	}

	return v.String(), nil
}
