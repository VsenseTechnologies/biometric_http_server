package routers

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"vsensetech.in/go_fingerprint_server/initilize"
)


func InitRouter(db *sql.DB , router *mux.Router , rdb *redis.Client , ctx context.Context){
	router.HandleFunc("/{id}/init", initilize.NewInitInstance(db , rdb , ctx).InitilizeTables)
}