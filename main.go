package main

import (
	"log"
	"net/http"
	"os"
	"sync"

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
		log.Println("Unable to load the Environment Variable please check and try again...")
	}
	
	//Creating Database Connection
	databaseConnection := database.DatabaseConnection{
		DatabaseName: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}
	db , err := databaseConnection.ConnectToDatabase()
	if err != nil {
		log.Println("Unable To Connect To Database Please Check...")
	}
	defer db.Close()
	
	//Declaring Router and Mutex
	router  := mux.NewRouter()
	mut := &sync.Mutex{}
	
	//Routes
	router.Use(middlewares.RouteMiddleware)
	routers.AuthRouter(db , mut, router)
	routers.InitRouter(db, router)
	routers.AdminRouters(db, mut, router)
	
	
	//Starting The Server
	port := os.Getenv("SERVER_PORT")
	log.Println("Server has Started and is running at PORT ",port)
	log.Fatal(http.ListenAndServe(port, router))
}