package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lpernett/godotenv"
	"vsensetech.in/go_fingerprint_server/database"
	"vsensetech.in/go_fingerprint_server/middlewares"
	"vsensetech.in/go_fingerprint_server/routers"
)

func main() {
	//Loading The Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	//Creating Database Connection
	databaseConnection := database.DatabaseConnection{
		DatabaseURL: os.Getenv("DB_URL"),
	}
	db, err := databaseConnection.ConnectToDatabase()
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()

	//Routes
	router.Use(middlewares.RouteMiddleware)
	router.Use(middlewares.JwtMiddleware)
	routers.AuthRouter(db, router)
	routers.InitRouter(db, router)
	routers.AdminRouters(db, router)
	routers.UserRoutes(db, router)
	routers.TimeRouters(db, router)

	//Starting The Server
	port := os.Getenv("SERVER_PORT")
	log.Println("Server has Started and is running at PORT ", port)
	log.Fatal(http.ListenAndServe(port, router))

}
