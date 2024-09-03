package controllers

import (
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type ManageUsersController struct{
	manageUserRepository models.ManageUsersRepository
}

func NewManageUsersController(mur models.ManageUsersRepository) *ManageUsersController {
	return &ManageUsersController{
		manageUserRepository: mur,
	}
}

func(muc *ManageUsersController) GiveUserAccessController(w http.ResponseWriter , r *http.Request){
	if err := muc.manageUserRepository.GiveUserAccess(&r.Body); err != nil{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
}