package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type StudentFingerprintDataController struct {
	studentFingerprintDataRepository models.StudentFingerprintDataRepository
	rdb *redis.Client
	ctx *context.Context
}

func NewStudentFingerprintDataController(m models.StudentFingerprintDataRepository , rdb *redis.Client , ctx *context.Context) *StudentFingerprintDataController {
	return &StudentFingerprintDataController{
		studentFingerprintDataRepository: m,
		rdb: rdb,
		ctx: ctx,
	}
}

func(sfdc *StudentFingerprintDataController) LoadDataController(w http.ResponseWriter , r * http.Request) {
	data , err := sfdc.studentFingerprintDataRepository.LoadData(&r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}

	err = sfdc.rdb.Set(*sfdc.ctx , "load" , data , 0).Err()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success" , Data: data})
}