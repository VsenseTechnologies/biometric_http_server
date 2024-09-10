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
		return fmt.Errorf("invalid credentials")
	}
	newStudent.StudentID = uuid.New().String()

	if _ , err := sfr.db.Exec("INSERT INTO fingerprintdata(student_id , unit_id , fingerprint) VALUES($1 , $2 , $3)",newStudent.StudentID , newStudent.UnitID , newStudent.FingerprintData); err != nil {
		return fmt.Errorf("unable to add the student fingerprint details")
	}
	
	var queryString = fmt.Sprintf("INSERT INTO %s(student_id , student_name , student_usn , department) VALUES($1 , $2 , $3 , $4)",newStudent.UnitID)
	if _ , err := sfr.db.Exec(queryString , newStudent.StudentID , newStudent.StudentName , newStudent.StudentUSN , newStudent.Department); err != nil {
		return err
	}
	return nil
}

func(sfr *StudentFingerprintRepo) FetchStudentDetails(reader *io.ReadCloser) ([]models.StudentDetailsModel , error){
	var unit models.FingerprintMachinesModel
	if err := json.NewDecoder(*reader).Decode(&unit); err != nil {
		return nil,fmt.Errorf("invalid unit")
	}
	var queryString = fmt.Sprintf("SELECT student_id , student_name , student_usn FROM %s", unit.UnitID)
	res , err := sfr.db.Query(queryString)
	if err != nil {
		return nil,fmt.Errorf("unable to fetch student details")
	}
	defer res.Close()

	var student models.StudentDetailsModel
	var students []models.StudentDetailsModel
	for res.Next(){
		if err := res.Scan(&student.StudentID,&student.StudentName , &student.StudentUSN); err != nil {
			return nil,fmt.Errorf("unable to add students")
		}
		students = append(students, student)
	}
	if res.Err() != nil {
		return nil,fmt.Errorf("something went wrong")
	}
	return students , nil
}

func(sfr *StudentFingerprintRepo) FetchStudentLogHistory(reader *io.ReadCloser) ([]models.StudentLogHistoryModel , error) {
	var student models.StudentDetailsModel
	if err := json.NewDecoder(*reader).Decode(&student); err != nil {
		return nil,fmt.Errorf("invalid studentID")
	}
	res , err := sfr.db.Query("SELECT login , logout , date FROM attendence WHERE student_id=$1",student.StudentID)
	if err != nil {
		return nil,fmt.Errorf("unable to fetch loghistory")
	}
	defer res.Close()

	var log models.StudentLogHistoryModel
	var logs []models.StudentLogHistoryModel
	for res.Next(){
		if err := res.Scan(&log.LoginTime , &log.LogoutTime , &log.Date); err != nil {
			return nil , fmt.Errorf("unable to fetch log  history")
		}
		logs = append(logs, log)
	}
	if res.Err() != nil {
		return nil ,fmt.Errorf("something went wrong")
	}
	return logs , nil
}

func(sfr *StudentFingerprintRepo) DeleteStudent(reader *io.ReadCloser) error {
	var studentCred models.StudentOperationModel
	if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
		return fmt.Errorf("invalid studentcredentials")
	}
	var queryString = fmt.Sprintf("DELETE FROM %s WHERE student_id=$1" , studentCred.UnitID)
	if _ , err := sfr.db.Exec(queryString , studentCred.StudentID); err != nil {
		return fmt.Errorf("unable to delete student")
	}
	return nil
}

func(sfr *StudentFingerprintRepo) UpdateStudent(reader *io.ReadCloser) error {
	var studentCred models.StudentOperationModel
	if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	var queryString = fmt.Sprintf("UPDATE %s SET student_name=$1 , student_usn=$2 WHERE student_id=$3" , studentCred.UnitID) 
	if _ , err := sfr.db.Exec(queryString , studentCred.StudentName , studentCred.StudentUSN , studentCred.StudentID); err != nil {
		return fmt.Errorf("unable to update student")
	}
	return nil
}
