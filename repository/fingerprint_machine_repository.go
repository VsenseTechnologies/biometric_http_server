package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type FingerprintMachineRepo struct {
	db  *sql.DB
	mut *sync.Mutex
}

func NewFingerprintMachineRepo(db *sql.DB, mut *sync.Mutex) *FingerprintMachineRepo {
	return &FingerprintMachineRepo{
		db,
		mut,
	}
}

func (umr *FingerprintMachineRepo) FetchAllMachines(reader *io.ReadCloser) ([]models.FingerprintMachinesModel, error) {
	// Locking The Process To Prevent Crashing
	umr.mut.Lock()
	defer umr.mut.Unlock()

	//
	var userID string
	if err := json.NewDecoder(*reader).Decode(&userID); err != nil {
		return nil, fmt.Errorf("invalid credendials")
	}
	res, err := umr.db.Query("SELECT unit_id , online FROM biometric WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get id")
	}
	defer res.Close()

	var userMachines []models.FingerprintMachinesModel
	var userMachine models.FingerprintMachinesModel

	for res.Next() {
		err := res.Scan(&userMachine.UnitID, &userMachine.Status)
		if err != nil {
			return nil, fmt.Errorf("invalid credendials")
		}
		userMachines = append(userMachines, userMachine)
	}
	if res.Err() != nil {
		return nil, fmt.Errorf("something went wrong")
	}
	return userMachines, nil
}

func (umr *FingerprintMachineRepo) DeleteMachine(reader *io.ReadCloser) error {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var machine models.FingerprintMachinesModel
	if err := json.NewDecoder(*reader).Decode(&machine); err != nil {
		return fmt.Errorf("unable to process request")
	}
	query := fmt.Sprintf("DROP TABLE %s", machine.UnitID)

	if _, err := umr.db.Exec(query); err != nil {
    	return fmt.Errorf("unable to delete table")
	}
	if _, err := umr.db.Exec("DELETE FROM biometric WHERE unit_id=$1", machine.UnitID); err != nil {
		return fmt.Errorf("unable to delete unit")
	}
	return nil
}

func (umr *FingerprintMachineRepo) AddMachine(reader *io.ReadCloser) error {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var newMachine models.FingerprintMachinesModel

	if err := json.NewDecoder(*reader).Decode(&newMachine); err != nil {
		return nil
	}


	query := fmt.Sprintf("CREATE TABLE %s (student_id VARCHAR(100), student_name VARCHAR(50) NOT NULL, student_usn VARCHAR(20) NOT NULL, department VARCHAR(20) NOT NULL , FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id)  ON DELETE CASCADE)", newMachine.UnitID)

	if _, err := umr.db.Exec(query); err != nil {
    	return fmt.Errorf("unable to create table")
	}


	if _, err := umr.db.Exec("INSERT INTO biometric(user_id , unit_id , online) VALUES($1 , $2 , $3)", newMachine.UserID, newMachine.UnitID, false); err != nil {
		return nil
	}
	return nil
}
