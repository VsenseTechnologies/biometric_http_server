package initialize

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"vsensetech.in/go_fingerprint_server/payload"
)

type Init struct {
	db *sql.DB
}

func NewInitInstance(db *sql.DB) *Init {
	return &Init{
		db,
	}
}

// Function to handle table creation with error handling inside a transaction
func (i *Init) createTable(tx *sql.Tx, query string) error {
	_, err := tx.Exec(query)
	return err
}

// InitializeTables method to create all tables with atomicity (transaction)
func (i *Init) InitializeTables(w http.ResponseWriter, r *http.Request) {
	// Start a new transaction
	tx, err := i.db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: "Failed to begin transaction: " + err.Error()})
		return
	}

	// Table creation queries
	tables := []string{
		`CREATE TABLE IF NOT EXISTS admin (
			user_id VARCHAR(100) PRIMARY KEY, 
			user_name VARCHAR(50) NOT NULL, 
			password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(100) PRIMARY KEY, 
			user_name VARCHAR(50) NOT NULL, 
			password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS biometric (
			user_id VARCHAR(100), 
			unit_id VARCHAR(50) PRIMARY KEY, 
			online BOOLEAN NOT NULL, 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS fingerprintdata (
			student_id VARCHAR(100) PRIMARY KEY, 
			student_unit_id VARCHAR(100), 
			unit_id VARCHAR(50), 
			fingerprint VARCHAR(1000), 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS attendance (
			student_id VARCHAR(100), 
			student_unit_id VARCHAR(100), 
			unit_id VARCHAR(50), 
			date VARCHAR(20), 
			login VARCHAR(20), 
			logout VARCHAR(20), 
			FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE, 
			FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS times (
			user_id VARCHAR(200), 
			morning_start VARCHAR(20), 
			morning_end VARCHAR(20), 
			afternoon_start VARCHAR(20), 
			afternoon_end VARCHAR(20), 
			evening_start VARCHAR(20), 
			evening_end VARCHAR(20), 
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS inserts(
			unit_id VARCHAR(200),
			student_unit_id VARCHAR(200),
			fingerprint_data VARCHAR(1000)
		)`,
		`CREATE TABLE IF NOT EXISTS deletes(
			unit_id VARCHAR(200),
			student_unit_id VARCHAR(200)
		)`,
	}

	// Iterate over the table creation queries and execute them in a transaction
	for _, query := range tables {
		if err := i.createTable(tx, query); err != nil {
			tx.Rollback() // Rollback the transaction if any query fails
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: "Failed to create table: " + err.Error()})
			return
		}
	}

	// Commit the transaction if all queries are successful
	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: "Failed to commit transaction: " + err.Error()})
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Tables initialized successfully"})
}
