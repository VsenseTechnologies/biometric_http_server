package models

import "io"


type StudentDetails struct{
	StudentID string `json:"student_id"`
	StudentName string `json:"name"`
	StudentUSN string `json:"usn"`
	StudentDepartment string `json:"department"`
}

type StudentLogHistory struct {
	LoginTime string `json:"login_time"`
	LogoutTime string `json:"logout_time"`
}

type StudentDetailsRepository interface{
	FetchAllStudentDetails(*io.ReadCloser) ([]StudentDetails , error)
	DeleteStudent(*io.ReadCloser) (error)
	UpdateStudentDetails(*io.ReadCloser) (error)
	FetchStudentLogHistory(*io.ReadCloser)
}

type CreateDocument struct{
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

type CreateDocumentRepository interface{
	DownloadDetails(*io.ReadCloser) error
}