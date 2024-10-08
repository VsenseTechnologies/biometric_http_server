package models

import "io"

type AuthDetails struct {
	Name     string `json:"user_name"`
	Password string `json:"password"`
}

type AuthDetailsRepository interface {
	Register(*io.ReadCloser, string) (string, error)
	Login(*io.ReadCloser, string) (string, error)
}
