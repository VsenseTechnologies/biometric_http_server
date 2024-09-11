package routers

import (
	"context"
	"database/sql"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func FingerprintMachineRouters(db *sql.DB , mut *sync.Mutex , router *mux.Router , rdb *redis.Client , ctx *context.Context){ 
	repo := repository.NewStudentFingerprintDataRepo(db , mut)
	cont := controllers.NewStudentFingerprintDataController(repo , rdb , *ctx)

	router.HandleFunc("/users/load" , cont.LoadDataController).Methods("POST")
}