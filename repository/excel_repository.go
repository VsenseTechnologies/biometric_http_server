package repository

import (
	"database/sql"
	"fmt"
	"io"
	"sync"

	"github.com/xuri/excelize/v2"
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

func(ar *AttendenceRepo) CreateAttendenceSheet(reader *io.ReadCloser) error {
	file := excelize.NewFile()
	defer func(){
		if err := file.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	index , err := file.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	file.SetCellValue("Sheet1" , "A2" , "Hello")
	file.SetActiveSheet(index)
	if err := file.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
	return nil
}