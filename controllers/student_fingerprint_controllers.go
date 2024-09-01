package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)


type StudentFingerprintController struct{
	studentFingerprintRepository models.StudentFingerprintRepository
}

func NewStudentFingerprintController(studentFingerprintRepository models.StudentFingerprintRepository) *StudentFingerprintController {
	return &StudentFingerprintController{
		studentFingerprintRepository,
	}
}

func(sfc *StudentFingerprintController) RegisterStudentController(w http.ResponseWriter , r *http.Request){
	if err := sfc.studentFingerprintRepository.RegisterStudent(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errs := json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()}); errs != nil {
			log.Println(errs.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if errs := json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"}); errs != nil {
		log.Println(errs.Error())
	}
}

func(sfc *StudentFingerprintController) FetchStudentDetailsController(w http.ResponseWriter , r *http.Request){
	data , err := sfc.studentFingerprintRepository.FetchStudentDetails(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errs := json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()}); errs != nil {
			log.Println(errs.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if errs := json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success" , Data: data}); errs != nil {
		log.Println(errs.Error())
	}
}

func(sfc *StudentFingerprintController) FetchStudentLogHistory(w http.ResponseWriter , r *http.Request){
	data , err := sfc.studentFingerprintRepository.FetchStudentLogHistory(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errs := json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()}); errs != nil {
			log.Println(errs.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if errs := json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success" , Data: data}); errs != nil {
		log.Println(errs.Error())
	}
}

func(sfc *StudentFingerprintController) DeleteStudentController(w http.ResponseWriter , r *http.Request){
	if err := sfc.studentFingerprintRepository.DeleteStudent(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errs := json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()}); errs != nil {
			log.Println(errs.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if errs := json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"}); errs != nil {
		log.Println(errs.Error())
	}
}

func(sfc *StudentFingerprintController) UpdateStudentController(w http.ResponseWriter , r *http.Request){
	if err := sfc.studentFingerprintRepository.UpdateStudent(&r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errs := json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()}); errs != nil {
			log.Println(errs.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if errs := json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"}); errs != nil {
		log.Println(errs.Error())
	}
}