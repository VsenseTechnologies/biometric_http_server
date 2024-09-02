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
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var user models.UsersModel
	if err := json.NewDecoder(*reader).Decode(&user); err != nil {
		return nil, err
	}
	res, err := umr.db.Query("SELECT unit_id , online FROM biometric WHERE user_id=$1", user.UserID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var userMachines []models.FingerprintMachinesModel
	var userMachine models.FingerprintMachinesModel

	for res.Next() {
		err := res.Scan(&userMachine.UnitID, &userMachine.Status)
		if err != nil {
			return nil, err
		}
		userMachines = append(userMachines, userMachine)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	return userMachines, nil
}

func (umr *FingerprintMachineRepo) DeleteMachine(reader *io.ReadCloser) error {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var machine models.FingerprintMachinesModel
	if err := json.NewDecoder(*reader).Decode(&machine); err != nil {
		return err
	}
	query := fmt.Sprintf("DROP TABLE %s", machine.UnitID)

	if _, err := umr.db.Exec(query); err != nil {
    	return err
	}
	if _, err := umr.db.Exec("DELETE FROM biometric WHERE unit_id=$1", machine.UnitID); err != nil {
		return err
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


	query := fmt.Sprintf("CREATE TABLE %s (student_id VARCHAR(100), name VARCHAR(50) NOT NULL, usn VARCHAR(20) NOT NULL, department VARCHAR(20) NOT NULL , FOREIGN KEY student_id REFERENCES fingerprintdata(student_id) ON DELETE CASCADE)", newMachine.UnitID)

	if _, err := umr.db.Exec(query); err != nil {
    	return err
	}


	if _, err := umr.db.Exec("INSERT INTO biometric(user_id , unit_id , online) VALUES($1 , $2 , $3)", newMachine.UserID, newMachine.UnitID, false); err != nil {
		return nil
	}
	return nil
}
