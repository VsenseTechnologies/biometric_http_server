package models

import "io"

type StudentFingerprintRegistrationModel struct {
	StudentID       string `json:"student_id,omitempty"`
	StudentUnitID   string `json:"student_unit_id,omitempty"`
	StudentName     string `json:"student_name,omitempty"`
	StudentUSN      string `json:"student_usn,omitempty"`
	Department      string `json:"department,omitempty"`
	UnitID          string `json:"unit_id,omitempty"`
	FingerprintData string `json:"fingerprint_data,omitempty"`
}

type StudentDetailsModel struct {
	StudentID     string `json:"student_id,omitempty"`
	StudentUnitID string `json:"student_unit_id,omitempty"`
	StudentName   string `json:"student_name,omitempty"`
	StudentUSN    string `json:"student_usn,omitempty"`
	Department    string `json:"department,omitempty"`
}

type StudentLogHistoryModel struct {
	Date       string `json:"date,omitempty"`
	LoginTime  string `json:"login,omitempty"`
	LogoutTime string `json:"logout,omitempty"`
}
type StudentOperationModel struct {
	StudentID     string `json:"student_id,omitempty"`
	UnitID        string `json:"unit_id,omitempty"`
	StudentUnitID string `json:"student_unit_id,omitempty"`
	StudentName   string `json:"student_name,omitempty"`
	StudentUSN    string `json:"student_usn,omitempty"`
	Department    string `json:"department,omitempty"`
}

type StudentFingerprintRepository interface {
	RegisterStudent(*io.ReadCloser) error
	DeleteStudent(*io.ReadCloser) error
	UpdateStudent(*io.ReadCloser) error
	FetchStudentDetails(*io.ReadCloser) ([]StudentDetailsModel, error)
	FetchStudentLogHistory(*io.ReadCloser) ([]StudentLogHistoryModel, error)
}
