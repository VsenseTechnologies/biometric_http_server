package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func AdminRouters(db *sql.DB, router *mux.Router) {
	//Admin Users Operation
	usersRepo := repository.NewUsersRepo(db)
	userCont := controllers.NewUsersController(usersRepo)

	router.HandleFunc("/admin/getusers", userCont.FetchAllUsersController).Methods("POST")
	// End Users Operation

	//Admin User Management Operation
	manageUserRepo := repository.NewManageUserRepo(db)
	manageUserCont := controllers.NewManageUsersController(manageUserRepo)

	router.HandleFunc("/admin/giveaccess", manageUserCont.GiveUserAccessController).Methods("POST")
	//End User Management Operation

	// Admin User Machine Operations
	userMachinesRepo := repository.NewFingerprintMachineRepo(db)
	userMachineCont := controllers.NewFingerprintMachineController(userMachinesRepo)

	router.HandleFunc("/admin/getmachines", userMachineCont.FetchAllMachinesController).Methods("POST")
	router.HandleFunc("/admin/addmachine", userMachineCont.AddMachineController).Methods("POST")
	router.HandleFunc("/admin/deletemachine", userMachineCont.DeleteMachineController).Methods("POST")
	// End User Machine Operations
}
