package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type DatabaseConnection struct{
	DatabaseURL string 
}

func (dc *DatabaseConnection) ConnectToDatabase() (*sql.DB , error){
	db , err := sql.Open("postgres", dc.DatabaseURL)
	if err != nil {
		return nil , err
	}
	return db , nil
}

