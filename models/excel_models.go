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

type Times struct{
	MorningStart string `json:"morning_start"`
	MorningEnd string `json:"morning_end"`
	AfternoonStart string `json:"afternoon_start"`
	AfternoonEnd string `json:"afternoon_end"`
	EveningStart string `json:"evening_start"`
	EveningEnd string `json:"evening_end"`
}

type AttendenceRepository interface {
	CreateAttendenceSheet(*io.ReadCloser) (*excelize.File, error)
}
