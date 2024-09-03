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

type Auth struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewAuth(db *sql.DB , mut *sync.Mutex) *Auth {
	return &Auth{
		db,
		mut,
	}
}

func(a *Auth) Register(reader *io.ReadCloser , urlPath string) (string , error) {
	a.mut.Lock()
	defer a.mut.Unlock()
	//Creating a new variable of type AdminAuthDetails
	var newUser models.AuthDetails
	
	//Decoding the json from reader to the newly created variale
	if err := json.NewDecoder(*reader).Decode(&newUser); err != nil {
		return "",fmt.Errorf("please enter valid details")
	}
	
	//Hashing the password
	hashpass , err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return "",fmt.Errorf("somthing went wrong while hashing")
	}
	//Execuiting the query and Creating new UUID and returning error if present
	var newUID = uuid.New().String()
	if _ , err := a.db.Exec("INSERT INTO "+urlPath+"(user_id , user_name , password) VALUES($1 , $2 , $3)", newUID , newUser.Name , hashpass); err != nil {
		return "",fmt.Errorf("unable to create user")
	}
	
	
	//Creating JWT token and Setting Cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":newUID,
		"user_name":newUser.Name,
		"expiry": time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString , err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "",fmt.Errorf("failed to create auth token")
	}

	//Return JWT token if No error
	return tokenString,nil
}

func(a *Auth) Login(reader *io.ReadCloser , urlPath string)  (string , error) {
	a.mut.Lock()
	defer a.mut.Unlock()
	//Creating a new variable of type AdminAuthDetails
	var userIns models.AuthDetails
	var dbUser models.AuthDetails
	var UID string
	
	//Decoding the json from reader to the newly created variale
	if err := json.NewDecoder(*reader).Decode(&userIns); err !=  nil {
		return "",fmt.Errorf("please enter valid details")
	}
	
	//Querying User from Database
		err := a.db.QueryRow("SELECT user_id , user_name , password FROM "+urlPath+" WHERE user_name=$1", userIns.Name).Scan(&UID, &dbUser.Name , &dbUser.Password)
		if err != nil {
			return "",fmt.Errorf("user is invalid")
		}

	
	//Comparing HashedPassword with Normal Password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userIns.Password))
	if err != nil {
		return "",fmt.Errorf("failed to validate password")
	}
	
	//Creating JWT token and Setting Cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":UID,
		"user_name":dbUser.Name,
		"expiry": time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString , err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "",fmt.Errorf("failed to create auth token")
	}
	
	//JWT token if No Error
	return tokenString,nil
}