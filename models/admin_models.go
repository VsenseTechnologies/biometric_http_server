package models

import "io"

// Users Struct
type Users struct {
	UserName string `json:"username"`
	UserID   string `json:"user_id"`
}

// Machines Struct
type UserMachines struct {
	UnitID string `json:"unit_id"`
	Status bool   `json:"online"`
}
type UserNewMachine struct {
	UnitID string `json:"unit_id"`
	UserID string `json:"user_id"`
}

// Struct to Add New Users
type ManageUsers struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UsersRepository interface {
	FetchAllUsers() ([]Users, error)
}

type UserMachineRepository interface {
	FetchAllMachines(*io.ReadCloser) ([]UserMachines, error)
	DeleteMachine(*io.ReadCloser) error
	AddMachine(*io.ReadCloser) error
}

type ManageUsersRepository interface {
	GiveUserAccess(*io.ReadCloser) error
}
