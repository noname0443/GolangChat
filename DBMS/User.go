package DBMS

import (
	"GolangChat/Utility"
	"database/sql"
	"errors"
	"time"
)

type User struct{
	Id int
	Username string
	Password string
	SessionPass *string
	SessionTime *time.Time
}

type Session struct{
	SessionPass string
	SessionTime time.Time
}

func readUserData(rows *sql.Rows) (error, []User) {
	var users []User
	for rows.Next(){
		p := User{}
		err := rows.Scan(&p.Username, &p.Password, &p.SessionPass, &p.SessionTime, &p.Id)
		if err != nil{
			return err, nil
		}
		users = append(users, p)
	}
	return nil, users
}

// GetUserBySession checks that user is exists by session information if it's correct returns user information
func GetUserBySession(db *sql.DB, session_pass string) (error, User) {
	rows, err := db.Query(`SELECT * FROM users WHERE session_pass = $1;`, session_pass)
	if err != nil {
		return err, User{}
	}
	err, resultUsers := readUserData(rows)
	if err != nil {
		return err, User{}
	}
	if len(resultUsers) == 1 {
		if time.Since(*resultUsers[0].SessionTime) > (time.Hour * 24) {
			return errors.New("session is too old"), User{}
		} else {
			return nil, resultUsers[0]
		}
	} else {
		return errors.New("can't find user"), User{}
	}
}

// GetUser checks that user is exists and if it's correct returns user information
func GetUser(db *sql.DB, username string, password string) (error, User) {
	hash := Utility.GetMD5Hash(password)
	rows, err := db.Query(`SELECT * FROM users WHERE username = $1 AND password = $2;`, username, hash)
	if err != nil {
		return err, User{}
	}
	err, resultUsers := readUserData(rows)
	if err != nil {
		return err, User{}
	}
	if len(resultUsers) == 1 {
		return nil, resultUsers[0]
	} else {
		return errors.New("Can't find user"), User{}
	}
}

// UpdateSession if user exists then get it and update if it's necessary user session
func UpdateSession(db *sql.DB, username string, password string) (error, User) {
	err, user := GetUser(db, username, password)
	if err != nil {
		return err, User{}
	}
	if user.SessionPass == nil || user.SessionTime == nil {
		user.SessionTime = new(time.Time)
		user.SessionPass = new(string)
	}
	if time.Since(*user.SessionTime) > (time.Hour * 24) {
		*user.SessionTime = time.Now()
		hash := Utility.GetMD5Hash(user.SessionTime.String())
		*user.SessionPass = hash
		_, err := db.Exec(`UPDATE users SET session_pass=$1, session_time=$2 WHERE id = $3;`, *user.SessionPass, *user.SessionTime, user.Id)
		if err != nil {
			return err, User{}
		}
		return nil, user
	} else {
		return nil, user
	}
}

// RegisterUser registers user in service and then update user's session
func RegisterUser(db *sql.DB, username string, password string) (error, User) {
	hash := Utility.GetMD5Hash(password)
	_, err := db.Exec(`INSERT INTO public.users(username, password) VALUES ($1, $2);`, username, hash)
	if err != nil {
		return err, User{}
	}
	return UpdateSession(db, username, password)
}