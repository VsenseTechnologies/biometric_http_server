package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type StudentFingerprintDataRepo struct {
	db *sql.DB
	mut *sync.Mutex
}

func NewStudentFingerprintDataRepo(db *sql.DB , mut *sync.Mutex) *StudentFingerprintDataRepo {
	return &StudentFingerprintDataRepo{
		db,
		mut,
	}
}

func(sfdr *StudentFingerprintDataRepo) LoadData(reader *io.ReadCloser) ([]models.StudentFingerprintData , error){
	var reqSFDs []models.StudentFingerprintData
	var dbSFDs []models.StudentFingerprintData
	var dbSFD models.StudentFingerprintData

	if err := json.NewDecoder(*reader).Decode(&reqSFDs); err != nil {
		return nil , err
	}

	res , err := sfdr.db.Query("SELECT student_id , unit_id , fingerprint FROM fingerprintdata WHERE unit_id=$1",reqSFDs[1].UnitID)
	if err != nil {
		return nil,err
	}
	defer res.Close()

	for res.Next(){
		if err := res.Scan(&dbSFD.StudentID , &dbSFD.UnitID); err != nil {
			return nil , err
		}
		for _ , id := range reqSFDs {
			if dbSFD.StudentID != id.StudentID {
				dbSFDs = append(dbSFDs , dbSFD)
				break
			}
		}
	} 
	if res.Err() != nil {
		return nil , res.Err()
	}

	return dbSFDs , nil
}