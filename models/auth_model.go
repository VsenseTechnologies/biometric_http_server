package models

import "io"

type AuthDetails struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type AuthRepository interface {
	Register(*io.ReadCloser, string) (string, error)
	Login(*io.ReadCloser, string) (string, error)
}
