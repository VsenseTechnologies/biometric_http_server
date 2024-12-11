package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/controllers"
	"vsensetech.in/go_fingerprint_server/repository"
)


func AuthRouter(db *sql.DB , router *mux.Router){
	repo := repository.NewAuth(db)
	cont := controllers.NewAuthController(repo)
	
	router.HandleFunc("/{id}/register", cont.RegisterController).Methods("POST")
	router.HandleFunc("/{id}/login", cont.LoginController).Methods("POST")
}