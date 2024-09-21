package models

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type AttendenceModel struct {
	StudentName string `json:"student_name"`
	StudentUSN  string `json:"student_usn"`
	Date        string `json:"date"`
	Login       string `json:"login"`
	Logout      string `json:"logout"`
}

type AttendenceRepository interface {
	CreateAttendenceSheet(*io.ReadCloser) (*excelize.File,error)
}
