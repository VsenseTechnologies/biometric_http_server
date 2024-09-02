package models

type UsersModel struct {
	UserName string `json:"user_name"`
	UserID   string `json:"user_id"`
	UserPassword string `json:"password"`
}

type UsersRepository interface {
	FetchAllUsers() ([]UsersModel, error)
}