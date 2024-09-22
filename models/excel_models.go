package models

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type AttendenceLogs struct {
	Date string `json:"date"`
	Login string `json:"login"`
	Logout string `json:"logout"`
}

type AttendenceStudent struct{
	StudentID string `json:"student_id"`
	StudentName string `json:"student_name"`
	StudentUSN string `json:"student_usn"`
	StudentUnitID string `json:"student_unit_id"`
	UnitID string `json:"unit_id"`
}

type AttendenceRepository interface {
	CreateAttendenceSheet(*io.ReadCloser) (*excelize.File, error)
}
