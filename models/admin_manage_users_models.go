package models

import "io"

type ManageUsers struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type ManageUsersRepository interface {
	GiveUserAccess(*io.ReadCloser) error
}
