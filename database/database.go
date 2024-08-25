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
	db , err := sql.Open("postgres", "postgresql://fingerprint_user:VAYbmJfOZYyJDBDKL1U7BRxv6OoqRR1h@dpg-cr4r595umphs73drro10-a/fingerprint")
	if err != nil {
		return nil , err
	}
	return db , nil
}

