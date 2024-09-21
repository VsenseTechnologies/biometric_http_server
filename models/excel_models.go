package models

import "io"

type AttendenceModel struct {
	Date   string `json:"date"`
	Login  string `json:"login"`
	Logout string `json:"logout"`
}

type AttendenceRepository interface {
	CreateAttendenceSheet(*io.ReadCloser) error
}