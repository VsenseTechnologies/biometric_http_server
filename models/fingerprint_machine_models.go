package models

import "io"

type FingerprintMachinesModel struct {
	UnitID string `json:"unit_id"`
	UserID string `json:"user_id"`
	Status bool   `json:"online"`
}

type FingerprintMachinesRepository interface {
	FetchAllMachines(*io.ReadCloser) ([]FingerprintMachinesModel, error)
	DeleteMachine(*io.ReadCloser) error
	AddMachine(*io.ReadCloser) error
}