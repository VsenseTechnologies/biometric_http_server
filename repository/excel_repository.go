package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"vsensetech.in/go_fingerprint_server/models"
)

type AttendenceRepo struct {
	db  *sql.DB
}

func NewAttendenceRepo(db *sql.DB) *AttendenceRepo {
	return &AttendenceRepo{
		db,
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
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#000000FF"}, // Light green background
			Pattern: 1,                   // Solid pattern
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

	// Set date headers for the attendance sheet starting from column "C"
	startDate := details.StartDate
	endDate := details.EndDate
	err = setAttendanceDateHeaders(file, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Fetch students
	students, err := FetchStudents(ar.db, details.UnitID)
	if err != nil {
		return nil, err
	}

	// Mark attendance and update Excel file
	update, err := MarkAttendance(ar.db, file, students, startDate, endDate, details.UserID)
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(index)
	return update, nil
}

func setAttendanceDateHeaders(file *excelize.File, startDate string, endDate string) error {
	// Parse the start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return err
	}

	// Create a style for bold text
	style, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#000000FF"}, // Light green background
			Pattern: 1,                   // Solid pattern
		},
	})
	if err != nil {
		return err
	}

	// Iterate through each date between startDate and endDate
	col := 0 // Starting from column 0, which will map to 'C'
	for currentDate := start; !currentDate.After(end); currentDate = currentDate.AddDate(0, 0, 1) {
		day := currentDate.Format("02") // Get only the day part

		column := columnIndexToLetter(col + 2) // Offset by 2 because columns start at "C"
		cell := column + "1"

		// Set the date value and apply the bold style
		file.SetCellValue("Sheet1", cell, day)
		file.SetColWidth("Sheet1", column, column, 5) // Set column width to 5
		file.SetCellStyle("Sheet1", cell, cell, style) // Apply bold style

		col++
	}
	return nil
}

// Fetch times from the 'times' table
func FetchTimes(db *sql.DB, userId string) (models.Times, error) {
	var times models.Times
	query := `SELECT morning_start, morning_end, afternoon_start, afternoon_end, evening_start, evening_end FROM times WHERE user_id=$1`
	err := db.QueryRow(query, userId).Scan(&times.MorningStart, &times.MorningEnd, &times.AfternoonStart, &times.AfternoonEnd, &times.EveningStart, &times.EveningEnd)
	if err != nil {
		return times, err
	}
	fmt.Println(times)
	return times, nil
}

// Mark attendance in the Excel file with the new logic
func MarkAttendance(db *sql.DB, file *excelize.File, data []models.AttendenceStudent, startDate string, endDate string, userId string) (*excelize.File, error) {
	l := 2 // Starting row for attendance entries

	// Fetch the times for attendance periods
	times, err := FetchTimes(db, userId)
	if err != nil {
		return nil, err
	}

	for _, student := range data {
		query := `SELECT date, login, logout FROM attendance WHERE student_id = $1 AND date::date BETWEEN $2::date AND $3::date`
		res, err := db.Query(query, student.StudentID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		defer res.Close()

		file.SetCellValue("Sheet1", "A"+strconv.Itoa(l), student.StudentName)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(l), student.StudentUSN)

		attendanceMap := make(map[string]string)

		for res.Next() {
			var log models.AttendenceLogs
			if err := res.Scan(&log.Date, &log.Login, &log.Logout); err != nil {
				return nil, err
			}

			// Determine attendance based on login/logout and the defined time ranges
			status := determineAttendance(log.Login, log.Logout, times)
			attendanceMap[log.Date] = status
		}

		// Mark each day with the corresponding attendance status
		// Mark each day with the corresponding attendance status
		// Parse the startDate and endDate as time.Time
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)

		// Mark each day with the corresponding attendance status
		currentDate := start
		for !currentDate.After(end) { // Ensure it includes the end date
			column := dateToColumn(currentDate.Format("2006-01-02"), startDate)
			if status, exists := attendanceMap[currentDate.Format("2006-01-02")]; exists {
				file.SetCellValue("Sheet1", column+strconv.Itoa(l), status)
			} else {
				file.SetCellValue("Sheet1", column+strconv.Itoa(l), "A") // Mark absent
			}
			currentDate = currentDate.AddDate(0, 0, 1) // Move to the next day
		}

		l++ // Move to the next student
	}
	return file, nil
}
func normalizeTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

// Determine attendance status based on login and logout times
func determineAttendance(login, logout string, times models.Times) string {
	// Parse the times from the Times struct with full hour:minute:second format
	morningEnd, _ := time.Parse("15:04:05", times.MorningEnd)
	afternoonStart, _ := time.Parse("15:04:05", times.AfternoonStart)
	afternoonEnd, _ := time.Parse("15:04:05", times.AfternoonEnd)
	eveningStart, _ := time.Parse("15:04:05", times.EveningStart)
	eveningEnd, _ := time.Parse("15:04:05", times.EveningEnd)

	// Parse the login and logout times with full hour:minute:second format
	loginTime, _ := time.Parse("15:04:05", login)
	logoutTime, _ := time.Parse("15:04:05", logout)

	// Normalize times to ignore the date component (comparing only time-of-day)
	loginTime = normalizeTime(loginTime)
	logoutTime = normalizeTime(logoutTime)
	morningEnd = normalizeTime(morningEnd)
	afternoonStart = normalizeTime(afternoonStart)
	afternoonEnd = normalizeTime(afternoonEnd)
	eveningStart = normalizeTime(eveningStart)
	eveningEnd = normalizeTime(eveningEnd)

	// Check if the student was present for the full day (P)
	if loginTime.Before(morningEnd) && logoutTime.After(eveningStart) && logoutTime.Before(eveningEnd) {
		return "P" // Full-day Present
	}

	// Check for Morning Present (MP)
	if loginTime.Before(morningEnd) && logoutTime.After(afternoonStart) && logoutTime.Before(afternoonEnd) {
		return "MP" // Morning Present
	}

	// Check for Afternoon Present (AP)
	if loginTime.After(afternoonStart) && logoutTime.Before(eveningEnd) {
		return "AP" // Afternoon Present
	}

	// If none of the above conditions match, mark as NC (No Conflict)
	return "NC" // Not Present or No Valid Attendance
}

// Helper function to map dates to Excel columns based on the startDate
func dateToColumn(date string, startDate string) string {
	baseDate, _ := time.Parse("2006-01-02", startDate)
	targetDate, _ := time.Parse("2006-01-02", date)
	diff := int(targetDate.Sub(baseDate).Hours() / 24)
	return columnIndexToLetter(2 + diff) // 'C' is the 3rd column, hence offset by 2
}

// Helper function to calculate Excel column letter (e.g., 0 -> "A", 25 -> "Z", 26 -> "AA", etc.)
func columnIndexToLetter(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
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
