package utils

import (
	"net/url"
)

func CheckURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return err
	}
	url, err := url.Parse(urlStr)
	if err != nil || url.Scheme == "" || url.Host == "" {
		return err
	}
	return nil
}
