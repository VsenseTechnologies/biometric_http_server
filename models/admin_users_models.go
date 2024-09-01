package models

type UsersModel struct {
	UserName string `json:"username"`
	UserID   string `json:"user_id"`
	UserPassword string `json:"password"`
}

type UsersRepository interface {
	FetchAllUsers() ([]UsersModel, error)
}