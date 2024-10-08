package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	_ "net/http/pprof"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/lpernett/godotenv"
	"vsensetech.in/go_fingerprint_server/database"
	"vsensetech.in/go_fingerprint_server/middlewares"
	"vsensetech.in/go_fingerprint_server/routers"
)

func main(){
	//Loading The Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	
	//Creating Database Connection
	databaseConnection := database.DatabaseConnection{
		DatabaseURL: os.Getenv("DB_URL"),
	}
	db , err := databaseConnection.ConnectToDatabase()
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()

	opt , err := redis.ParseURL("redis://default:dv5AD7fFG1jLNacDyJVYSjAfwlgNQwlX@redis-16241.c305.ap-south-1-1.ec2.redns.redis-cloud.com:16241")

	if err != nil {
		log.Println(err.Error())
	}
	//Redis Connection
	rdb := redis.NewClient(opt)

	ctx := context.Background()
	
	//Declaring Router and Mutex
	router  := mux.NewRouter()
	mut := &sync.Mutex{}
	
	//Routes
	router.Use(middlewares.RouteMiddleware)
	router.Use(middlewares.JwtMiddleware)
	routers.AuthRouter(db , mut, router)
	routers.InitRouter(db, router , rdb , ctx)
	routers.AdminRouters(db, mut, router , rdb , ctx)
	routers.UserRoutes(db , mut , router , rdb , ctx)
	routers.FingerprintMachineRouters(db , mut , router , rdb , &ctx)
	
	
	//Starting The Server
	port := os.Getenv("SERVER_PORT")
	log.Println("Server has Started and is running at PORT ",port)
	log.Fatal(http.ListenAndServe(port, router))
	
}