package models

import "io"

type MachineDetailsModel struct{
	UserID string `json:"user_id"`
	UnitID string `json:"unit_id"`
	Status string `json:"online"`
}

type MachineDetailsRepository interface{
	FetchAllMachines(*io.ReadCloser) ([]MachineDetailsModel , error)
}