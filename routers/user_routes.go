package routers

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func UserRoutes(db *sql.DB , mut *sync.Mutex , router *mux.Router){
	repo := repository.NewFingerprintMachineRepo(db , mut)
	cont := controllers.NewFingerprintMachineController(repo)

	http.HandleFunc("/users/getmachines" , cont.FetchAllMachinesController)

	repos := repository.NewStudentFingerprintRepo(db , mut)
	conts := controllers.NewStudentFingerprintController(repos)

	http.HandleFunc("/users/register" , conts.RegisterStudentController)
	http.HandleFunc("/users/fetchstudents" , conts.FetchStudentDetailsController)
	http.HandleFunc("/users/studentslog" , conts.FetchStudentLogHistory)
	http.HandleFunc("/users/deletestudent",conts.DeleteStudentController)
	http.HandleFunc("/users/updatestudent" , conts.UpdateStudentController)
}