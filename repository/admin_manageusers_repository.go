package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/smtp"
	"os"

	"golang.org/x/crypto/bcrypt"
	"vsensetech.in/go_fingerprint_server/models"
)

type ManageUserRepo struct{
	db *sql.DB
}

func NewManageUserRepo(db *sql.DB) *ManageUserRepo {
	return &ManageUserRepo{
		db,
	}
}


func (mur *ManageUserRepo) GiveUserAccess(reader *io.ReadCloser) error {
	var newUser models.ManageUsers
	var password string
	if err := json.NewDecoder(*reader).Decode(&newUser); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	if err := mur.db.QueryRow("SELECT password FROM users WHERE user_name=$1", newUser.UserName).Scan(&password); err != nil {
		return fmt.Errorf("invalid password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(newUser.Password)); err != nil {
		return fmt.Errorf("unable to validate password")
	}

	// Define email subject and body
	mailSubject := "Access Granted to VSENSE Fingerprint Software"
	mailBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
			<style>
				body { font-family: Arial, sans-serif; color: #000000; background-color: #ffffff; margin: 0; padding: 0; }
				.container { width: 90%%; margin: 0 auto; padding: 20px; }
				.header { border-bottom: 2px solid #000000; padding-bottom: 10px; margin-bottom: 20px; }
				.header h2 { margin: 0; color: #000000; }
				.content { margin-bottom: 20px; }
				.footer { border-top: 2px solid #000000; padding-top: 10px; text-align: center; color: #000000; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>%s</h2>
				</div>
				<div class="content">
					<p>Hello <b>%s</b>,</p>
					<p>We are pleased to inform you that access has been granted to the VSENSE Fingerprint Software.</p>
					<p><b>Username:</b> %s</p>
					<p><b>Password:</b> %s</p>
					<p>Please ensure to keep your credentials secure and do not share them with unauthorized individuals.</p>
					<p>If you need any assistance, feel free to contact our support team.</p>
				</div>
				<div class="footer">
					<p>Best regards,<br>The VSENSE Team</p>
				</div>
			</div>
		</body>
		</html>
	`, mailSubject, mailSubject, newUser.UserName, newUser.UserName, newUser.Password)

	// Convert the mail body to bytes
	mailMessage := []byte(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", mailSubject, mailBody))

	// Set up the SMTP client
	mailDetails := smtp.PlainAuth(
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
		mailMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
