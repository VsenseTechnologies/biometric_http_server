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
    ctx context.Context
}

func NewStudentFingerprintDataController(m models.StudentFingerprintDataRepository, rdb *redis.Client, ctx context.Context) *StudentFingerprintDataController {
    return &StudentFingerprintDataController{
        studentFingerprintDataRepository: m,
        rdb: rdb,
        ctx: ctx,
    }
}

func (sfdc *StudentFingerprintDataController) LoadDataController(w http.ResponseWriter, r *http.Request) {
    // Load data from repository
    data, err := sfdc.studentFingerprintDataRepository.LoadData(&r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
        return
    }

    // Create a Redis list key
    redisKey := "load"

    // Push each item in the data list into Redis
    for _, item := range data {
        // Marshal individual item to JSON
        jsonData, err := json.Marshal(item)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
            return
        }

        // Push the JSON data to Redis list
        _, err = sfdc.rdb.RPush(sfdc.ctx, redisKey, jsonData).Result()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
            return
        }
    }

    // Respond with success and include the data
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: data})
}
