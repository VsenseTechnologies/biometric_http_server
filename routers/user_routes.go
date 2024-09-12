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

func UserRoutes(db *sql.DB , mut *sync.Mutex , router *mux.Router , rdb *redis.Client , ctx context.Context){
	repo := repository.NewFingerprintMachineRepo(db , mut)
	cont := controllers.NewFingerprintMachineController(repo)

	router.HandleFunc("/users/getmachines" , cont.FetchAllMachinesController).Methods("POST")

	repos := repository.NewStudentFingerprintRepo(db , mut , rdb , ctx)
	conts := controllers.NewStudentFingerprintController(repos)

	router.HandleFunc("/users/registerstudent" , conts.RegisterStudentController).Methods("POST")
	router.HandleFunc("/users/fetchstudents" , conts.FetchStudentDetailsController).Methods("POST")
	router.HandleFunc("/users/studentslog" , conts.FetchStudentLogHistory).Methods("POST")
	router.HandleFunc("/users/deletestudent",conts.DeleteStudentController).Methods("POST")
	router.HandleFunc("/users/updatestudent" , conts.UpdateStudentController).Methods("POST")
}
