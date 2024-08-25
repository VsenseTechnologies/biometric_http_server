package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type UsersController struct{
	userRepository models.UsersRepository
}

func NewUsersController(ur models.UsersRepository) *UsersController {
	return &UsersController{
		userRepository: ur,
	}
}

func(uc *UsersController) FetchAllUsersController(w http.ResponseWriter , r *http.Request){
	data , err := uc.userRepository.FetchAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success" , Data: data})
}