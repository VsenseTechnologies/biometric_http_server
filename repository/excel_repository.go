package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/xuri/excelize/v2"
	"vsensetech.in/go_fingerprint_server/models"
)

type AttendenceRepo struct {
	db  *sql.DB
	mut *sync.Mutex
}

func NewAttendenceRepo(db *sql.DB, mut *sync.Mutex) *AttendenceRepo {
	return &AttendenceRepo{
		db,
		mut,
	}
}

func (ar *AttendenceRepo) CreateAttendenceSheet(reader *io.ReadCloser) (*excelize.File, error) {
	file := excelize.NewFile()
	var details models.AttendenceStudent
	if err := json.NewDecoder(*reader).Decode(&details); err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	index, err := file.NewSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	style, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true, 
			Size: 14, 
		},
	})
	if err != nil {
		return nil, err
	}

	file.SetCellStyle("Sheet1", "A1", "A1", style)
	file.SetColWidth("Sheet1", "A", "A", 30)
	file.SetCellValue("Sheet1", "A1", "Name")
	file.SetCellStyle("Sheet1", "B1", "B1", style)
	file.SetColWidth("Sheet1", "B", "B", 30)
	file.SetCellValue("Sheet1", "B1", "USN")

	for j := 67; j < 90; j++ {
		file.SetColWidth("Sheet1", string(j), string(j), 5)
	}
	var students []models.AttendenceStudent
	students , err = FetchStudents(ar.db , details.UnitID)
	if err != nil {
		return nil , err
	}
	file , err = MarkAttendance(ar.db , file , students , "2024-10-01" , "2024-10-29")
	if err != nil {
		return nil , err
	}
	file.SetActiveSheet(index)
	return file, nil
}

func FetchStudents(db *sql.DB , unitId string) ([]models.AttendenceStudent,error) {
	query := fmt.Sprintf("SELECT student_id , student_name , student_usn , student_unit_id FROM %s" , unitId)
	res , err := db.Query(query)
	if err != nil {
		return nil , err
	}
	defer res.Close()
	var students []models.AttendenceStudent
	for res.Next() {
		var student models.AttendenceStudent
		if err := res.Scan(&student.StudentID , &student.StudentName , &student.StudentUSN , &student.StudentUnitID); err != nil {
			return nil , err
		}
		students = append(students, student)
	}
	if res.Err() != nil {
		return nil , err
	}
	return students , nil
}

func MarkAttendance(db *sql.DB , file *excelize.File , data []models.AttendenceStudent , startDate string , endDate string) (*excelize.File,error){
	l := 2
	for _ , j := range data {
		query := fmt.Sprintf("SELECT date , login , logout FROM attendence WHERE student_id=$1 AND date::date BETWEEN %s ::date AND %s ::date",startDate , endDate)
		res , err := db.Query(query,j.StudentID)
		if err != nil {
			return nil , err
		}
		defer res.Close()
		for res.Next(){
			var log models.AttendenceLogs
			if err := res.Scan(&log.Date , &log.Login , &log.Logout); err != nil {
				return nil,err
			}
			for k := 67 ; k <= 90 ; k++ {
				file.SetCellValue("Sheet1" , string(k)+strconv.Itoa(l) , "P")
			}
		}
		if res.Err() != nil {
			return nil , err
		}
		l++;
	}
	return file , nil
}
 