package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
	"vsensetech.in/go_fingerprint_server/models"
)

type StudentFingerprintRepo struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewStudentFingerprintRepo(db *sql.DB , mut *sync.Mutex) *StudentFingerprintRepo {
	return &StudentFingerprintRepo{
		db,
		mut,
	}
}

func(sfr *StudentFingerprintRepo) RegisterStudent(reader *io.ReadCloser) error {
	var newStudent models.StudentFingerprintRegistrationModel
	if err := json.NewDecoder(*reader).Decode(&newStudent); err != nil {
		return err
	}
	newStudent.StudentID = uuid.New().String()
	var queryString = fmt.Sprintf("INSERT INTO %s(student_id , student_name , student_usn , student_department) VALUES($1 , $2 , $3 , $4)",newStudent.UnitID)
	if _ , err := sfr.db.Exec(queryString , newStudent.StudentID , newStudent.StudentName , newStudent.StudentUSN , newStudent.Department); err != nil {
		return err
	}
	if _ , err := sfr.db.Exec("INSERT INTO fingerprint_details(student_id , unit_id , fingerprint_data) VALUES($1 , $2 , $3)",newStudent.StudentID , newStudent.UnitID , newStudent.FingerprintData); err != nil {
		return err
	}
	return nil
}

func(sfr *StudentFingerprintRepo) FetchStudentDetails(reader *io.ReadCloser) ([]models.StudentDetailsModel , error){
	var unitId string
	if err := json.NewDecoder(*reader).Decode(&unitId); err != nil {
		return nil,err
	}
	var queryString = fmt.Sprintf("SELECT student_name , student_usn FROM %s", unitId)
	res , err := sfr.db.Query(queryString)
	if err != nil {
		return nil,err
	}
	defer res.Close()

	var student models.StudentDetailsModel
	var students []models.StudentDetailsModel
	for res.Next(){
		if err := res.Scan(&student.StudentName , &student.StudentUSN); err != nil {
			return nil,err
		}
		students = append(students, student)
	}
	if res.Err() != nil {
		return nil,res.Err()
	}
	return students , nil
}

func(sfr *StudentFingerprintRepo) FetchStudentLogHistory(reader *io.ReadCloser) ([]models.StudentLogHistoryModel , error) {
	var studentID string
	if err := json.NewDecoder(*reader).Decode(&studentID); err != nil {
		return nil,err
	}
	res , err := sfr.db.Query("SELECT login , logout , date FROM attendence WHERE unit_id=$1",studentID)
	if err != nil {
		return nil,err
	}
	defer res.Close()

	var log models.StudentLogHistoryModel
	var logs []models.StudentLogHistoryModel
	for res.Next(){
		if err := res.Scan(&log.LoginTime , &log.LogoutTime , &log.Date); err != nil {
			return nil , err
		}
		logs = append(logs, log)
	}
	if res.Err() != nil {
		return nil , res.Err()
	}
	return logs , nil
}

func(sfr *StudentFingerprintRepo) DeleteStudent(reader *io.ReadCloser) error {
	var studentCred models.StudentOperationModel
	if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
		return err
	}
	var queryString = fmt.Sprintf("DELETE FROM %s WHERE student_id=$1" , studentCred.UnitID)
	if _ , err := sfr.db.Exec(queryString , studentCred.StudentID); err != nil {
		return err
	}
	return nil
}

func(sfr *StudentFingerprintRepo) UpdateStudent(reader *io.ReadCloser) error {
	var studentCred models.StudentOperationModel
	if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
		return err
	}
	var queryString = fmt.Sprintf("UPDATE %s SET student_name=$1 , student_usn=$2 WHERE student_id=$3" , studentCred.UnitID) 
	if _ , err := sfr.db.Exec(queryString , studentCred.StudentName , studentCred.StudentUSN , studentCred.StudentID); err != nil {
		return err
	}
	return nil
}