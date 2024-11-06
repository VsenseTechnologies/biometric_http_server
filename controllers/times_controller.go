package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type timeController struct {
	timeRepository models.TimeInterface
}

func NewTimeController(timeRepository models.TimeInterface) *timeController{
	return &timeController{
		timeRepository,
	}
}

func(tc *timeController) SetTimeController(w http.ResponseWriter , r *http.Request){
	if err := tc.timeRepository.SetTime(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}