package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type UserMachineController struct {
	userMachineRepository models.UserMachineRepository
}

func NewUserMachineController(umr models.UserMachineRepository) *UserMachineController {
	return &UserMachineController{
		userMachineRepository: umr,
	}
}

func (umc *UserMachineController) FetchAllMachinesController(w http.ResponseWriter, r *http.Request) {
	data, err := umc.userMachineRepository.FetchAllMachines(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: data})
	return
}

func (umc *UserMachineController) DeleteMachineController(w http.ResponseWriter, r *http.Request) {
	if err := umc.userMachineRepository.DeleteMachine(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
	return
}

func (umc *UserMachineController) AddMachineController(w http.ResponseWriter, r *http.Request) {
	if err := umc.userMachineRepository.AddMachine(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
	return
}
