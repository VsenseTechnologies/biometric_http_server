package routers

import (
	"database/sql"

	"github.com/gorilla/mux"
	initialize "vsensetech.in/go_fingerprint_server/initilize"
)


func InitRouter(db *sql.DB , router *mux.Router){
	init := initialize.NewInitInstance(db)
	router.HandleFunc("/users/init", init.InitializeTables)
}
