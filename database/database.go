package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type DatabaseConnection struct{
	Username string
	Password string
	DatabaseName string
	SSLMode string
}

func (dc *DatabaseConnection) ConnectToDatabase() (*sql.DB , error){
	db , err := sql.Open("postgres", "user="+dc.Username+" "+"password="+dc.Password+" "+"dbname="+dc.DatabaseName+" "+"sslmode="+dc.SSLMode)
	if err != nil {
		return nil , err
	}
	return db , nil
}

