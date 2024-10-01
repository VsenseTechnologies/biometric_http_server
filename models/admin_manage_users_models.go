package models

import "io"

type ManageUsers struct {
	Email    string `json:"email,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Password string `json:"password,omitempty"`
}

type ManageUsersRepository interface {
	GiveUserAccess(*io.ReadCloser) error
}
