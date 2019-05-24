package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InsertToDB(filename string, mistakes int, notKnow int, timeSpent int64) error {
	db, err := sql.Open("sqlite3", "./report.db")
	if err != nil {
		return err
	}
	defer db.Close()

	err = createTableIfNotExist(db)
	if err != nil {
		return err
	}

	sqlStmt := "INSERT INTO REPORTS(date, filename, mistakes, not_know, time_spent) " +
		"VALUES('%s', '%s', %d, %d, %d);"

	_, err = db.Exec(fmt.Sprintf(sqlStmt, now(), filename, mistakes, notKnow, timeSpent))
	if err != nil {
		return err
	}

	return nil
}

func createTableIfNotExist(db *sql.DB) error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS reports (
		        date     TEXT,
				filename TEXT,
				mistakes INTEGER,
				not_know INTEGER,
				time_spent INTEGER
			   );
			   `

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func now() string {
	now := time.Now()

	return strconv.Itoa(now.Year()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Day())
}
