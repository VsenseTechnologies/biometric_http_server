package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type FingerprintMachineController struct {
	userMachineRepository models.FingerprintMachinesRepository
}

func NewFingerprintMachineController(umr models.FingerprintMachinesRepository) *FingerprintMachineController {
	return &FingerprintMachineController{
		userMachineRepository: umr,
	}
}

func (umc *FingerprintMachineController) FetchAllMachinesController(w http.ResponseWriter, r *http.Request) {
	data, err := umc.userMachineRepository.FetchAllMachines(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: data})
}

func (umc *FingerprintMachineController) DeleteMachineController(w http.ResponseWriter, r *http.Request) {
	if err := umc.userMachineRepository.DeleteMachine(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}

func (umc *FingerprintMachineController) AddMachineController(w http.ResponseWriter, r *http.Request) {
	if err := umc.userMachineRepository.AddMachine(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}
