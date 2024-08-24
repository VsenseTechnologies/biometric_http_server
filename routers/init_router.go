package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/initilize"
)


func InitRouter(db *sql.DB , router *mux.Router){
	router.HandleFunc("/{id}/init", initilize.NewInitInstance(db).InitilizeTables)
}