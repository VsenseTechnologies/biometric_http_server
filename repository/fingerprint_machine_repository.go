package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/go-redis/redis/v8"
	"vsensetech.in/go_fingerprint_server/models"
)

type FingerprintMachineRepo struct {
	db  *sql.DB
	mut *sync.Mutex
	rdb *redis.Client
	ctx context.Context
}

func NewFingerprintMachineRepo(db *sql.DB, mut *sync.Mutex , rdb *redis.Client , ctx context.Context) *FingerprintMachineRepo {
	return &FingerprintMachineRepo{
		db,
		mut,
		rdb,
		ctx,
	}
}

func (umr *FingerprintMachineRepo) FetchAllMachines(reader *io.ReadCloser) ([]models.FingerprintMachinesModel, error) {
	// Locking The Process To Prevent Crashing
	umr.mut.Lock()
	defer umr.mut.Unlock()

	var userCred models.UsersModel

	if err := json.NewDecoder(*reader).Decode(&userCred); err != nil {
		return nil, fmt.Errorf("invalid credendials")
	}
	res, err := umr.db.Query("SELECT unit_id , online FROM biometric WHERE user_id=$1", userCred.UserID)
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

	if _ , err := umr.rdb.Do(umr.ctx,"JSON.DEL" , "deletes" , "$."+machine.UnitID).Result(); err != nil {
		return err
	}

	if _ , err := umr.rdb.Do(umr.ctx,"JSON.DEL" , "inserts" , "$."+machine.UnitID).Result(); err != nil {
		return err
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

	
	if _ , err := umr.rdb.Do(umr.ctx,"JSON.SET" , "inserts" , "$."+newMachine.UnitID , "[]").Result(); err != nil {
		return err
	}

	if _ , err := umr.rdb.Do(umr.ctx,"JSON.SET" , "deletes" , "$."+newMachine.UnitID , "[]").Result(); err != nil {
		return err
	}


	query := fmt.Sprintf("CREATE TABLE %s (student_id VARCHAR(100) , student_unit_id VARCHAR(100) NOT NULL, student_name VARCHAR(50) NOT NULL, student_usn VARCHAR(20) NOT NULL, department VARCHAR(20) NOT NULL , FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id)  ON DELETE CASCADE)", newMachine.UnitID)

	if _, err := umr.db.Exec(query); err != nil {
    	return fmt.Errorf("unable to create table")
	}


	if _, err := umr.db.Exec("INSERT INTO biometric(user_id , unit_id , online) VALUES($1 , $2 , $3)", newMachine.UserID, newMachine.UnitID, false); err != nil {
		return err
	}


	return nil
}
