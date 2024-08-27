package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type UserMachineRepo struct {
	db  *sql.DB
	mut *sync.Mutex
}

func NewUserMachineRepo(db *sql.DB, mut *sync.Mutex) *UserMachineRepo {
	return &UserMachineRepo{
		db,
		mut,
	}
}

func (umr *UserMachineRepo) FetchAllMachines(reader *io.ReadCloser) ([]models.UserMachines, error) {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var user models.Users
	if err := json.NewDecoder(*reader).Decode(&user); err != nil {
		return nil, err
	}
	res, err := umr.db.Query("SELECT unit_id , online FROM biometric WHERE user_id=$1", user.UserID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var userMachines []models.UserMachines
	var userMachine models.UserMachines

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

func (umr *UserMachineRepo) DeleteMachine(reader *io.ReadCloser) error {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var machine models.UserMachines
	if err := json.NewDecoder(*reader).Decode(&machine); err != nil {
		return err
	}
	if _, err := umr.db.Exec("DELETE FROM biometric WHERE unit_id=$1", machine.UnitID); err != nil {
		return err
	}
	return nil
}

func (umr *UserMachineRepo) AddMachine(reader *io.ReadCloser) error {
	umr.mut.Lock()
	defer umr.mut.Unlock()
	var newMachine models.UserNewMachine

	if err := json.NewDecoder(*reader).Decode(&newMachine); err != nil {
		return nil
	}

	if _, err := umr.db.Exec("INSERT INTO biometric(user_id , unit_id , online) VALUES($1 , $2 , $3)", newMachine.UserID, newMachine.UnitID, false); err != nil {
		return nil
	}
	return nil
}
