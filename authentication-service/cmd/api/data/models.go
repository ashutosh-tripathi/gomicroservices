package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	UserID      string    `json:"user"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Status      bool      `json:"status"`
	Created_at  time.Time `json:"created_at"`
	Modified_at time.Time `json:"modified_at"`
}

const dbTimeout = time.Second * 3

// var Db *sql.DB

func checkDBConnection(db *sql.DB) {
	err := db.Ping()
	fmt.Println("error while pinging", err)
}
func (u *User) AddUser(db *sql.DB) (string, error) {
	checkDBConnection(db)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := "Insert into userdata.user_data values($1, $2,$3,$4,$5,$6)"

	if db == nil {
		fmt.Println("error creating database")
	}

	_, err := db.ExecContext(ctx, stmt, u.UserID, u.Password, u.Email, u.Status, u.Created_at, u.Modified_at)

	if err != nil {
		fmt.Println("error", err.Error())
		return "", err
	}
	return u.UserID, nil
}
func (u *User) GetUser(userid string, db *sql.DB) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := "select * from userdata.user_data where  user_id = $1"
	rows, err := db.QueryContext(ctx, stmt, userid)
	if err != nil {
		return User{}, err
	}
	var user User

	for rows.Next() {
		err := rows.Scan(
			&user.UserID,
			&user.Password,
			&user.Email,
			&user.Status,
			&user.Created_at,
			&user.Modified_at,
		)
		if err != nil {
			return User{}, err
		} else {
			break
		}

	}

	return user, nil

}
func (u User) GetAll(db *sql.DB) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []User
	stmt := "select * from userdata.user_data"
	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserID,
			&user.Password,
			&user.Email,
			&user.Status,
			&user.Created_at,
			&user.Modified_at,
		)
		if err != nil {
			return users, err
		} else {
			users = append(users, user)
		}
	}
	fmt.Printf("fetched user is %#v", users)
	return users, nil
}
func (u User) UpdateUser(user User, db *sql.DB) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `Update userdata.user_data set
	       password=$1,
		   email=$2,status=$3,
		   created_at=$4,updated_at=$5 where user_id=$6`
	res, err := db.ExecContext(ctx, stmt, user.Password, user.Email, user.Status, user.Created_at, user.Modified_at, u.UserID)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return rows, nil
}

func (u User) DeleteUser(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := "Delete from userdata.user_data where user_id=$1"
	_, err := db.ExecContext(ctx, stmt, u.UserID)
	if err != nil {
		return err
	} else {
		return nil
	}

}
func (u User) UserExists(db *sql.DB, userid string, pwd string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := "Select * from userdata.user_data where user_id=$1 and password=$2"
	rows, err := db.QueryContext(ctx, stmt, userid, pwd)

	if err != nil {
		fmt.Println("in eror found")
		if err == sql.ErrNoRows {
			fmt.Println("no rows found")
			return false, err
		}
		return false, err
	} else {
		if !rows.Next() {
			return false, nil
		} else {
			return true, nil
		}

	}

}
