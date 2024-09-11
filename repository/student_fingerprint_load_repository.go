package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	fmt.Println(reqSFDs)

	res , err := sfdr.db.Query("SELECT student_id , unit_id , fingerprint FROM fingerprintdata WHERE unit_id=$1",reqSFDs[0].UnitID)
	if err != nil {
		return nil,err
	}
	defer res.Close()

	for res.Next(){
		if err := res.Scan(&dbSFD.StudentID , &dbSFD.UnitID , &dbSFD.FingerprintData); err != nil {
			return nil , err
		}
		dbSFDs = append(dbSFDs, dbSFD)
	} 
	fmt.Println(dbSFDs)
	if res.Err() != nil {
		return nil , res.Err()
	}

	for _ , id := range reqSFDs{
		for i , sid := range dbSFDs{
			if id.StudentID == sid.StudentID {
				dbSFDs = removeElement(dbSFDs , i)
			}
		}
	}

	fmt.Println(reqSFDs)
	fmt.Println(dbSFDs)

	return dbSFDs , nil
}

func removeElement(slice []models.StudentFingerprintData, index int) []models.StudentFingerprintData {
    if index < 0 || index >= len(slice) {
        return slice // Index out of range
    }
    return append(slice[:index], slice[index+1:]...)
}