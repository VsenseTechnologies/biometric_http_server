package models

import "io"

type StudentFingerprintData struct {
	StudentID       string `json:"student_id,omitempty"`
	StudentUnitID   string `json:"student_unit_id,omitempty"`
	UnitID          string `json:"unit_id,omitempty"`
	FingerprintData string `json:"fingerprint,omitempty"`
}

type StudentFingerprintDataRepository interface {
	LoadData(*io.ReadCloser) ([]StudentFingerprintData, error)
}
