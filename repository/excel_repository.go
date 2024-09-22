package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

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

	// Set header style
	style, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})
	if err != nil {
		return nil, err
	}

	// Set column headers and styles
	file.SetCellStyle("Sheet1", "A1", "A1", style)
	file.SetColWidth("Sheet1", "A", "A", 30)
	file.SetCellValue("Sheet1", "A1", "Name")
	file.SetCellStyle("Sheet1", "B1", "B1", style)
	file.SetColWidth("Sheet1", "B", "B", 30)
	file.SetCellValue("Sheet1", "B1", "USN")

	// Set column widths for attendance dates
	for j := 67; j < 90; j++ {
		file.SetColWidth("Sheet1", string(j), string(j), 5)
	}

	// Fetch students
	students, err := FetchStudents(ar.db, details.UnitID)
	if err != nil {
		return nil, err
	}

	// Mark attendance and update Excel file
	update, err := MarkAttendance(ar.db, file, students, "2024-10-01", "2024-10-29")
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(index)
	return update, nil
}

func FetchStudents(db *sql.DB, unitId string) ([]models.AttendenceStudent, error) {
	// Use parameterized query to prevent SQL injection
	query := fmt.Sprintf("SELECT student_id, student_name, student_usn, student_unit_id FROM %s", unitId)
	res, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var students []models.AttendenceStudent
	for res.Next() {
		var student models.AttendenceStudent
		if err := res.Scan(&student.StudentID, &student.StudentName, &student.StudentUSN, &student.StudentUnitID); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	// Check for errors in result processing
	if err = res.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func MarkAttendance(db *sql.DB, file *excelize.File, data []models.AttendenceStudent, startDate string, endDate string) (*excelize.File, error) {
	l := 2 // Starting row for attendance entries
	for _, student := range data {
		// Parameterized query for fetching attendance logs
		query := `SELECT date, login, logout FROM attendence WHERE student_id = $1 AND date BETWEEN $2::date AND $3::date`
		res, err := db.Query(query, student.StudentID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		defer res.Close()

		// Set student info (name, USN)
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(l), student.StudentName)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(l), student.StudentUSN)

		// Iterate over attendance logs for this student
		for res.Next() {
			var log models.AttendenceLogs
			if err := res.Scan(&log.Date, &log.Login, &log.Logout); err != nil {
				return nil, err
			}

			// Convert log.Date to column letter (start from 'C')
			// Example logic assumes attendance dates are mapped to columns starting from 'C'
			column := dateToColumn(log.Date) // You need to implement this function
			file.SetCellValue("Sheet1", column+strconv.Itoa(l), "P") // Mark as present
		}

		// Check for errors during iteration
		if err := res.Err(); err != nil {
			return nil, err
		}

		l++ // Move to the next row for the next student
	}
	return file, nil
}

// Helper function to map dates to Excel columns (e.g., 2024-10-01 to 'C', 2024-10-02 to 'D', etc.)
func dateToColumn(date string) string {
	// Logic to calculate which column to use based on the date (e.g., "2024-10-01" = "C", "2024-10-02" = "D", etc.)
	// You need to map each date to a corresponding Excel column letter.
	// This is a placeholder logic:
	// Convert the date to an integer and map it to the correct column
	baseDate, _ := time.Parse("2006-01-02", "2024-10-01")
	targetDate, _ := time.Parse("2006-01-02", date)
	diff := int(targetDate.Sub(baseDate).Hours() / 24)
	return string(rune(67 + diff)) // 'C' is 67 in ASCII
}
