package models

import "io"

type FingerprintMachinesModel struct {
	UnitID string `json:"unit_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Status bool   `json:"online,omitempty"`
}

type FingerprintMachinesRepository interface {
	FetchAllMachines(*io.ReadCloser) ([]FingerprintMachinesModel, error)
	DeleteMachine(*io.ReadCloser) error
	AddMachine(*io.ReadCloser) error
}