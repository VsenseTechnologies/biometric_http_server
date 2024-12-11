package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func UserRoutes(db *sql.DB , router *mux.Router){
	repo := repository.NewFingerprintMachineRepo(db)
	cont := controllers.NewFingerprintMachineController(repo)

	router.HandleFunc("/users/getmachines" , cont.FetchAllMachinesController).Methods("POST")

	repos := repository.NewStudentFingerprintRepo(db)
	conts := controllers.NewStudentFingerprintController(repos)

	router.HandleFunc("/users/registerstudent" , conts.RegisterStudentController).Methods("POST")
	router.HandleFunc("/users/fetchstudents" , conts.FetchStudentDetailsController).Methods("POST")
	router.HandleFunc("/users/studentslog" , conts.FetchStudentLogHistory).Methods("POST")
	router.HandleFunc("/users/deletestudent",conts.DeleteStudentController).Methods("POST")
	router.HandleFunc("/users/updatestudent" , conts.UpdateStudentController).Methods("POST")
	reposit := repository.NewAttendenceRepo(db)
	control := controllers.NewAttendenceController(reposit)
	router.HandleFunc("/users/download" , control.CreateAttendenceSheetController).Methods("POST")
}
