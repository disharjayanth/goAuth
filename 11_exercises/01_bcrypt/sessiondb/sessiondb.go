package sessiondb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Session struct {
	Name string
	SID  string
}

var db *sql.DB
var err error

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

func (s *Session) Store() (bool, error) {
	result, err := db.Exec("INSERT INTO sessionid (name, sid) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = ? , sid = ? ", s.Name, s.SID, s.Name, s.SID)
	if err != nil {
		return false, fmt.Errorf("error while inserting sid into sessionid table: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error while getting sid's id after inserting: %w", err)
	}

	if count != 0 {
		return true, nil
	}

	return false, fmt.Errorf("unexpected error while inserting")
}

func Get(sid string) (string, error) {
	var session Session
	row := db.QueryRow("SELECT * FROM sessionid WHERE sid = ?", sid)
	if err := row.Scan(&session.Name, &session.SID); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with sid %s: %w", sid, err)
		}
		return "", fmt.Errorf("error while retreiving sid user")
	}

	return session.Name, nil
}
