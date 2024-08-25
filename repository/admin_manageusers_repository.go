package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/smtp"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"vsensetech.in/go_fingerprint_server/models"
)

type ManageUserRepo struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewManageUserRepo(db *sql.DB , mut *sync.Mutex) *ManageUserRepo {
	return &ManageUserRepo{
		db,
		mut,
	}
}

func(mur *ManageUserRepo) GiveUserAccess(reader *io.ReadCloser) error {
	var newUser models.ManageUsers
	var password string
	if err := json.NewDecoder(*reader).Decode(&newUser); err != nil {
		return err
	}
	if err := mur.db.QueryRow("SELECT password FROM users WHERE username=$1", newUser.UserName).Scan(&password); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(newUser.Password)); err != nil {
		return err
	}
	var mailSubject string = "Access to Fingerprint Software"
	var mailBody string = "This Mail consist of username and password for Accessing VSENSE Fingerprint Software\nUsername: "+newUser.UserName+"\npassword: "+newUser.Password
	var mailMessage string = mailSubject + "\n" + mailBody
	var mailDetails = smtp.PlainAuth(
		"", 
		os.Getenv("SMTP_USERNAME"), 
		os.Getenv("SMTP_PASSWORD"), 
		os.Getenv("SMTP_SERVICE"),
	)
	err := smtp.SendMail(
		"smtp.gmail.com:587", 
		mailDetails, 
		os.Getenv("SMTP_USERNAME"), 
		[]string{newUser.Email}, 
		[]byte(mailMessage),
	)
	if err != nil {
		return nil
	}
	return nil
}