package routers

import (
	"database/sql"
	"sync"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func UserRoutes(db *sql.DB , mut *sync.Mutex , router *mux.Router){
	repo := repository.NewFingerprintMachineRepo(db , mut)
	cont := controllers.NewFingerprintMachineController(repo)

	router.HandleFunc("/users/getmachines" , cont.FetchAllMachinesController)

	repos := repository.NewStudentFingerprintRepo(db , mut)
	conts := controllers.NewStudentFingerprintController(repos)

	router.HandleFunc("/users/register" , conts.RegisterStudentController)
	router.HandleFunc("/users/fetchstudents" , conts.FetchStudentDetailsController)
	router.HandleFunc("/users/studentslog" , conts.FetchStudentLogHistory)
	router.HandleFunc("/users/deletestudent",conts.DeleteStudentController)
	router.HandleFunc("/users/updatestudent" , conts.UpdateStudentController)
}