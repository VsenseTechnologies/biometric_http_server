package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"vsensetech.in/go_fingerprint_server/models"
)

type AuthDetailsRepo struct {
	db  *sql.DB
	mut *sync.Mutex
}

func NewAuth(db *sql.DB, mut *sync.Mutex) *AuthDetailsRepo {
	return &AuthDetailsRepo{
		db,
		mut,
	}
}

func (a *AuthDetailsRepo) Register(reader *io.ReadCloser, urlPath string) (string, error) {
	// Locking The Process To Avoid Crashes
	a.mut.Lock()
	defer a.mut.Unlock()
	// Creating New User Model to Store json to go Objects
	var newUser models.AuthDetails

	// Reading the json from Reader and Storing it on to Objects
	if err := json.NewDecoder(*reader).Decode(&newUser); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Hashing the Password for Security
	hashpass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("somthing went wrong")
	}
	// Creating a New User ID for the User
	var newUID = uuid.New().String()

	// Creating a JWT Request Verificaation Token to Validate Users Request
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        newUID,
		"user_name": newUser.Name,
		"expiry":    time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "", fmt.Errorf("failed to create auth token")
	}

	// Entering User Details on to Database
	if _, err := a.db.Exec("INSERT INTO "+urlPath+"(user_id , user_name , password) VALUES($1 , $2 , $3)", newUID, newUser.Name, hashpass); err != nil {
		return "", fmt.Errorf("unable to create user")
	}

	// Returning Generated Token to User
	return tokenString, nil
}

func (a *AuthDetailsRepo) Login(reader *io.ReadCloser, urlPath string) (string, error) {
	// Locking The Process To Avoid Crashes
	a.mut.Lock()
	defer a.mut.Unlock()

	// The reqData stores The Data Sent By User and The dbData Stores The Data Fetched From Database With Respect To User Data
	var reqData models.AuthDetails
	var dbData models.AuthDetails
	var UID string

	// Reading the json from Reader and Storing it on to Objects
	if err := json.NewDecoder(*reader).Decode(&reqData); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Querying User from Database
	err := a.db.QueryRow("SELECT user_id , user_name , password FROM "+urlPath+" WHERE user_name=$1", reqData.Name).Scan(&UID, &dbData.Name, &dbData.Password)
	if err != nil {
		return "", fmt.Errorf("unable to retrive the user data")
	}

	// Comparing HashedPassword with Normal Password
	err = bcrypt.CompareHashAndPassword([]byte(dbData.Password), []byte(reqData.Password))
	if err != nil {
		return "", fmt.Errorf("unable to validate password")
	}

	// Creating a JWT Request Verificaation Token to Validate Users Request
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        UID,
		"user_name": dbData.Name,
		"expiry":    time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "", fmt.Errorf("unable to create auth token")
	}

	// Sending The JWT Token to User
	return tokenString, nil
}
