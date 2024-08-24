package routers

import (
	"database/sql"
	"sync"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)


func AdminAuthRouter(db *sql.DB , mut *sync.Mutex , router *mux.Router){
	repo := repository.NewAuth(db, mut)
	cont := controllers.NewAuthController(repo)
	
	router.HandleFunc("/{id}/register", cont.RegisterController).Methods("POST")
	router.HandleFunc("/{id}/login", cont.LoginController).Methods("POST")
}