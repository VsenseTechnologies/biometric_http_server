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
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment;filename=Book1.xlsx")
	w.Header().Set("File-Name", "Book1.xlsx")
	f , err := ac.attendenceRepository.CreateAttendenceSheet(&r.Body); 
	if err != nil {
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
	}
	err = f.Write(w)
	if err != nil {
		http.Error(w, "Unable to generate Excel file", http.StatusInternalServerError)
		
	}
}