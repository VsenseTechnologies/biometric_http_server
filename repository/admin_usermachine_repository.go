package repository

import (
	"database/sql"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type UserMachineRepo struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewUserMachineRepo(db *sql.DB , mut *sync.Mutex) *UserMachineRepo {
	return &UserMachineRepo{
		db,
		mut,
	}
}

func(umr *UserMachineRepo) FetchAllMachines() ([]models.UserMachines , error ){
	res , err := umr.db.Query("SELECT unit_id , online FROM biometric")
	if err != nil {
		return nil,err
	}
	defer res.Close()
	
	var userMach
	return nil , nil
}