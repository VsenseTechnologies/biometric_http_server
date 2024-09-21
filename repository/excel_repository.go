package repository

import (
	"database/sql"
	"fmt"
	"io"
	"sync"

	"github.com/xuri/excelize/v2"
	"vsensetech.in/go_fingerprint_server/models"
)

type AttendenceRepo struct {
	db *sql.DB
	mut *sync.Mutex
}

func NewAttendenceRepo(db *sql.DB , mut *sync.Mutex) *AttendenceRepo {
	return &AttendenceRepo{
		db,
		mut,
	}
}

func(ar *AttendenceRepo) CreateAttendenceSheet(reader *io.ReadCloser) (*excelize.File , error) {
	file := excelize.NewFile()
	defer func(){
		if err := file.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	index , err := file.NewSheet("Sheet1")
	if err != nil {
		return nil,err
	}
	style, err := file.NewStyle(&excelize.Style{
        Font: &excelize.Font{
            Bold: true,        // Set text to bold
            Size: 14,   // Set font size to 14
        },
    })
	if err != nil {
		return nil,err
	}

	file.SetCellStyle("Sheet1","A1","A1",style)
	file.SetColWidth("Sheet1" , "A" , "A" , 30)
	file.SetCellStyle("Sheet1","B1","B1",style)
	file.SetColWidth("Sheet1" , "B" , "B" , 30)
	file.SetCellValue("Sheet1" , "A1" , "Name")
	file.SetCellValue("Sheet1" , "B1" , "USN")

	res , err := ar.db.Query("SELECT student_name , student_usn FROM vs24al001")
	if err != nil {
		return nil,err
	}

	defer res.Close()
	var student models.AttendenceModel
	i := 2
	for res.Next(){
		if err := res.Scan(&student.StudentName , &student.StudentUSN); err != nil {
			return nil,err
		}
		file.SetCellValue("Sheet1" , "A"+string(i) , student.StudentName)
		file.SetCellValue("Sheet1" , "B"+string(i) , student.StudentUSN)
		i++;
	}


	file.SetActiveSheet(index)
	if err := file.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
	return file,nil
}