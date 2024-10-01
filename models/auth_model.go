package models

import "io"

type AuthDetails struct {
	Name     string `json:"user_name,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthDetailsRepository interface {
	Register(*io.ReadCloser, string) (string, error)
	Login(*io.ReadCloser, string) (string, error)
}
