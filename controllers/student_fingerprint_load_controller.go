package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type StudentFingerprintDataController struct {
	studentFingerprintDataRepository models.StudentFingerprintDataRepository
}

func NewStudentFingerprintDataController(m models.StudentFingerprintDataRepository) *StudentFingerprintDataController {
	return &StudentFingerprintDataController{
		studentFingerprintDataRepository: m,
	}
}

func(sfdc *StudentFingerprintDataController) LoadDataController(w http.ResponseWriter , r * http.Request) {
	data , err := sfdc.studentFingerprintDataRepository.LoadData(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success" , Data: data})
}