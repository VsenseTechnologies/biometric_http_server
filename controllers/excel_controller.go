package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type AttendenceController struct {
	attendenceRepository models.AttendenceRepository
}

func NewAttendenceController(attendenceRepository models.AttendenceRepository) *AttendenceController {
	return &AttendenceController{
		attendenceRepository,
	}
}

func(ac *AttendenceController) CreateAttendenceSheetController(w http.ResponseWriter , r *http.Request){
	if err := ac.attendenceRepository.CreateAttendenceSheet(&r.Body); err != nil {
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
	}
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}