package models

import "io"

type StudentFingerprintRegistrationModel struct {
	StudentID       string `json:"student_id"`
	StudentUnitID   string `json:"student_unit_id"`
	StudentName     string `json:"student_name"`
	StudentUSN      string `json:"student_usn"`
	Department      string `json:"department"`
	UnitID          string `json:"unit_id"`
	FingerprintData string    `json:"fingerprint_data"`
}

type StudentDetailsModel struct {
	StudentID     string `json:"student_id"`
	StudentUnitID string `json:"student_unit_id"`
	StudentName   string `json:"student_name"`
	StudentUSN    string `json:"student_usn"`
}

type StudentLogHistoryModel struct {
	Date       string `json:"date"`
	LoginTime  string `json:"login"`
	LogoutTime string `json:"logout"`
}
type StudentOperationModel struct {
	StudentID   string `json:"student_id"`
	UnitID      string `json:"unit_id"`
	StudentUnitID string `json:"student_unit_id"`
	StudentName string `json:"student_name"`
	StudentUSN  string `json:"student_usn"`
}

type StudentFingerprintRepository interface {
	RegisterStudent(*io.ReadCloser) error
	DeleteStudent(*io.ReadCloser) error
	UpdateStudent(*io.ReadCloser) error
	FetchStudentDetails(*io.ReadCloser) ([]StudentDetailsModel, error)
	FetchStudentLogHistory(*io.ReadCloser) ([]StudentLogHistoryModel, error)
}
