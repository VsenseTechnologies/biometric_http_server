package initilize

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"vsensetech.in/go_fingerprint_server/payload"
)

type Init struct{
	db *sql.DB
}

func NewInitInstance(db *sql.DB) *Init{
	return &Init{
		db,
	}
}

func(i *Init) InitilizeTables(w http.ResponseWriter , r *http.Request){
	if _ , err := i.db.Exec("CREATE TABLE admin(user_id VARCHAR(100) PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(100) NOT NULL)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE users(user_id VARCHAR(100) PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(100) NOT NULL)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE biometric(user_id VARCHAR(100), unit_id VARCHAR(50) PRIMARY KEY , online BOOLEAN NOT NULL, FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE fingerprintdata(user_id VARCHAR(100), unit_id VARCHAR(50) PRIMARY KEY , online BOOLEAN NOT NULL, FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE, FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	if _ , err := i.db.Exec("CREATE TABLE attendence(user_id VARCHAR(100), unit_id VARCHAR(50), login TIMESTAMP DEFAULT CURRENT_TIMESTAMP, logout TIMESTAMP, FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE, FOREIGN KEY (unit_id) REFERENCES biometric(unit_id) ON DELETE CASCADE)"); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload.SimepleSuccessPayload{Message: "Success"})
}