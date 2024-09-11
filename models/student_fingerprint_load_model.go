package models

import "io"

type StudentFingerprintData struct {
	StudentID       string `json:"student_id"`
	UnitID          string `json:"unit_id"`
	FingerprintData string `json:"fingerprint"`
}

type StudentFingerprintDataRepository interface {
	LoadData(*io.ReadCloser) ([]StudentFingerprintData , error)
}