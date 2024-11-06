package routers

import (
	"database/sql"
	"sync"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)

func TimeRouters(db *sql.DB , mut *sync.Mutex , router *mux.Router){ 
	repo := repository.NewTimeRepository(db , mut)
	cont := controllers.NewTimeController(repo)

	router.HandleFunc("/users/time" , cont.SetTimeController)
}