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

func AdminRouters(db *sql.DB, mut *sync.Mutex, router *mux.Router , rdb *redis.Client , ctx context.Context) {
	//Admin Users Operation
	usersRepo := repository.NewUsersRepo(db, mut)
	userCont := controllers.NewUsersController(usersRepo)

	router.HandleFunc("/admin/getusers", userCont.FetchAllUsersController).Methods("POST")
	// End Users Operation

	//Admin User Management Operation
	manageUserRepo := repository.NewManageUserRepo(db, mut)
	manageUserCont := controllers.NewManageUsersController(manageUserRepo)

	router.HandleFunc("/admin/giveaccess", manageUserCont.GiveUserAccessController).Methods("POST")
	//End User Management Operation

	// Admin User Machine Operations
	userMachinesRepo := repository.NewFingerprintMachineRepo(db, mut , rdb , ctx)
	userMachineCont := controllers.NewFingerprintMachineController(userMachinesRepo)

	router.HandleFunc("/admin/getmachines", userMachineCont.FetchAllMachinesController).Methods("POST")
	router.HandleFunc("/admin/addmachine", userMachineCont.AddMachineController).Methods("POST")
	router.HandleFunc("/admin/deletemachine", userMachineCont.DeleteMachineController).Methods("POST")
	// End User Machine Operations
}
