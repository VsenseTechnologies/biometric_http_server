package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func TimeRouters(db *sql.DB , router *mux.Router){ 
	repo := repository.NewTimeRepository(db)
	cont := controllers.NewTimeController(repo)

	router.HandleFunc("/users/time" , cont.SetTimeController)
}