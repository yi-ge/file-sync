package utils

import (
	"errors"
	"net/url"
)

func CheckURL(urlStr string) error {
	obj, err := url.Parse(urlStr)

	if obj.Scheme == "" {
		return errors.New("missing scheme")
	}

	if obj.Host == "" {
		return errors.New("missing host")
	}

	if err != nil {
		return err
	}

	_, err = url.ParseRequestURI(urlStr)

	if err != nil {
		return err
	}

	return nil
}
