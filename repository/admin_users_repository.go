package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type UsersRepo struct{
	db *sql.DB
	mut *sync.Mutex
}

func NewUsersRepo(db *sql.DB , mut *sync.Mutex) *UsersRepo{
	return &UsersRepo{
		db,
		mut,
	}
}

func(ur *UsersRepo) FetchAllUsers() ([]models.UsersModel , error) {
	// Locking The Process to Avoid Crashes
	ur.mut.Lock()
	defer ur.mut.Unlock()

	// Getting All The Users Data From Database
	res , err := ur.db.Query("SELECT user_name , user_id FROM users")
	if err != nil {
		return nil ,fmt.Errorf("unable to fetch colleges")
	}
	defer res.Close()
	
	// Creating 2 go Models one is Slice which is of type UserModel
	var userList []models.UsersModel
	var user models.UsersModel
	
	// Going Through All the Rows And Creating go Objects and Storing it in the Slice of that model Type
	for res.Next() {
		err := res.Scan(&user.UserName , &user.UserID)
		if err != nil {
			return nil , fmt.Errorf("unable to add college")
		}
		userList = append(userList, user)
	}
	if res.Err() != nil {
		return nil , fmt.Errorf("something went wrong")
	}

	// Returning the list of Users Fetched From Database
	return userList , nil
}