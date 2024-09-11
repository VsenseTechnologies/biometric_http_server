package routers

import (
	"database/sql"
	"sync"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func FingerprintMachineRouters(db *sql.DB , mut *sync.Mutex , router *mux.Router){ 
	repo := repository.NewStudentFingerprintDataRepo(db , mut)
	cont := controllers.NewStudentFingerprintDataController(repo)

	router.HandleFunc("/users/load" , cont.LoadDataController).Methods("POST")
}