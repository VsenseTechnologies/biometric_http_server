package models

import "io"

type ManageUsers struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ManageUsersRepository interface {
	GiveUserAccess(*io.ReadCloser) error
}
