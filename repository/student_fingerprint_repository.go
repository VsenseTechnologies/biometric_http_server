package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"vsensetech.in/go_fingerprint_server/models"
)

type StudentFingerprintRepo struct{
	db *sql.DB
	mut *sync.Mutex
	rdb *redis.Client
	ctx context.Context
}

func NewStudentFingerprintRepo(db *sql.DB , mut *sync.Mutex , rdb *redis.Client , ctx context.Context) *StudentFingerprintRepo {
	return &StudentFingerprintRepo{
		db,
		mut,
		rdb,
		ctx,
	}
}

func(sfr *StudentFingerprintRepo) RegisterStudent(reader *io.ReadCloser) error {
	// Locking the process to prevent crashing
	sfr.mut.Lock()
	defer sfr.mut.Unlock()

	// Creating data models to store the data
	var newStudent models.StudentFingerprintRegistrationModel

	// Decoding the json data and storing it to the data model
	if err := json.NewDecoder(*reader).Decode(&newStudent); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	// Generating new unique id for student
	newStudent.StudentID = uuid.New().String()

	// Inserting fingerprint data to the Student fingerprint table
	if _ , err := sfr.db.Exec("INSERT INTO fingerprintdata(student_id , student_unit_id , unit_id , fingerprint) VALUES($1 , $2 , $3 , $4)",newStudent.StudentID, newStudent.StudentUnitID , newStudent.UnitID , newStudent.FingerprintData); err != nil {
		return err
	}
	
	// Inserting Student fingerprint data on to the machine table
	var queryString = fmt.Sprintf("INSERT INTO %s(student_id , student_unit_id , student_name , student_usn , department) VALUES($1 , $2 , $3 , $4 , $5)",newStudent.UnitID)
	if _ , err := sfr.db.Exec(queryString , newStudent.StudentID , newStudent.StudentUnitID , newStudent.StudentName , newStudent.StudentUSN , newStudent.Department); err != nil {
		return fmt.Errorf("unable to add data to machine")
	}


	return nil
}

func(sfr *StudentFingerprintRepo) FetchStudentDetails(reader *io.ReadCloser) ([]models.StudentDetailsModel , error){
	// Locking the process to prevent crashing
	sfr.mut.Lock()
	defer sfr.mut.Unlock()

	// Creating data models to store the data
	var unit models.FingerprintMachinesModel

	// Decoding the json data and storing on to the data model
	if err := json.NewDecoder(*reader).Decode(&unit); err != nil {
		return nil,fmt.Errorf("invalid unit")
	}

	// Querying the student details from the machine table
	var queryString = fmt.Sprintf("SELECT student_id , student_name , student_usn , student_unit_id FROM %s", unit.UnitID)
	res , err := sfr.db.Query(queryString)
	if err != nil {
		return nil,fmt.Errorf("unable to fetch student details")
	}
	defer res.Close()

	// Creating data models to store data
	var student models.StudentDetailsModel
	var students []models.StudentDetailsModel
	for res.Next(){
		if err := res.Scan(&student.StudentID,&student.StudentName , &student.StudentUSN , &student.StudentUnitID); err != nil {
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
	// Locking the process to prevent crashing
	sfr.mut.Lock()
	defer sfr.mut.Unlock()

	// Creating data model to store data
	var student models.StudentDetailsModel

	// Decoding the json data and storing on to the data model
	if err := json.NewDecoder(*reader).Decode(&student); err != nil {
		return nil,fmt.Errorf("invalid studentID")
	}

	// Executing the query to get data from database
	res , err := sfr.db.Query("SELECT login , logout , date FROM attendence WHERE student_id=$1",student.StudentID)
	if err != nil {
		return nil,fmt.Errorf("unable to fetch loghistory")
	}
	defer res.Close()

	// Creating data models for storing data
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

func (sfr *StudentFingerprintRepo) DeleteStudent(reader *io.ReadCloser) error {
    // Locking the process to prevent crashing
    sfr.mut.Lock()
    defer sfr.mut.Unlock()

    // Creating models to store data
    var studentCred models.StudentOperationModel

    // Decoding the JSON data and storing it into the model
    if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
        return fmt.Errorf("invalid student credentials: %w", err)
    }

    // Executing the delete query and deleting the data
    if _, err := sfr.db.Exec("DELETE FROM fingerprintdata WHERE student_id=$1", studentCred.StudentID); err != nil {
        return err
    }

	data , err := json.Marshal(map[string]string{
		"student_unit_id": studentCred.StudentUnitID,
	})
	if err != nil {
		return err
	}

	if _, err = sfr.rdb.Do(sfr.ctx, "JSON.ARRAPPEND", "deletes", "$."+studentCred.UnitID, data).Result(); err != nil {
		return err
	}
	
  
    return nil
}


func(sfr *StudentFingerprintRepo) UpdateStudent(reader *io.ReadCloser) error {
	// Locking the process to prevent crashing
	sfr.mut.Lock()
	defer sfr.mut.Unlock()

	// Creating student model to store the data
	var studentCred models.StudentOperationModel

	// Decoding the json data and storing on to the model
	if err := json.NewDecoder(*reader).Decode(&studentCred); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	// Querying the data and updating the details
	var queryString = fmt.Sprintf("UPDATE %s SET student_name=$1 , student_usn=$2 WHERE student_id=$3" , studentCred.UnitID) 
	if _ , err := sfr.db.Exec(queryString , studentCred.StudentName , studentCred.StudentUSN , studentCred.StudentID); err != nil {
		return fmt.Errorf("unable to update student")
	}
	return nil
}
