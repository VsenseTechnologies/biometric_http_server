package models

type UsersModel struct {
	UserName string `json:"user_name"`
	UserID   string `json:"user_id"`
}

type UsersRepository interface {
	FetchAllUsers() ([]UsersModel, error)
}