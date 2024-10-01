package models

type UsersModel struct {
	UserName string `json:"user_name,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

type UsersRepository interface {
	FetchAllUsers() ([]UsersModel, error)
}