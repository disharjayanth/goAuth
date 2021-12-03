package usersdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type User struct {
	ID       int
	Name     string
	Password string
}

func init() {
	// connection properties
	config := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		DBName:               "golangUsers",
		AllowNativePasswords: true,
	}

	// get a database handle
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

func (u *User) Store() (bool, error) {
	result, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", u.Name, u.Password)
	if err != nil {
		return false, fmt.Errorf("error adding user to user db: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, fmt.Errorf("error receiving user id: %w", err)
	}

	if id != 0 {
		return true, nil
	}

	return false, fmt.Errorf("unexpected error while inserting user")
}

func RetrieveUser(name string) (*User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE name = ?", name)
	if err := row.Scan(&user.ID, &user.Name, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with name: %s", name)
		}
		return nil, fmt.Errorf("error while retrieving user: %w", err)
	}

	return &user, nil
}
