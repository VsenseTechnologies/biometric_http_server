package models

import "io"

type TimeModel struct {
	UserID         string `json:"user_id"`
	MorningStart   string `json:"ms"`
	MorningEnd     string `json:"me"`
	AfternoonStart string `json:"as"`
	AfternoonEnd   string `json:"ae"`
	EveningStart   string `json:"es"`
	EveningEnd     string `json:"ee"`
}

type TimeInterface interface {
	SetTime(*io.ReadCloser) error
}
