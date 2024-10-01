package models

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type AttendenceLogs struct {
	Date string `json:"date,omitempty"`
	Login string `json:"login,omitempty"`
	Logout string `json:"logout,omitempty"`
}

type AttendenceStudent struct{
	StudentID string `json:"student_id,omitempty"`
	StudentName string `json:"student_name,omitempty"`
	StudentUSN string `json:"student_usn,omitempty"`
	StudentUnitID string `json:"student_unit_id,omitempty"`
	UnitID string `json:"unit_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
}

type Times struct{
	MorningStart string `json:"morning_start,omitempty"`
	MorningEnd string `json:"morning_end,omitempty"`
	AfternoonStart string `json:"afternoon_start,omitempty"`
	AfternoonEnd string `json:"afternoon_end,omitempty"`
	EveningStart string `json:"evening_start,omitempty"`
	EveningEnd string `json:"evening_end,omitempty"`
}

type AttendenceRepository interface {
	CreateAttendenceSheet(*io.ReadCloser) (*excelize.File, error)
}
