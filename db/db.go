package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbType = "sqlite3"
var dbFilename = "./report.db"

func Open() (*sql.DB, error) {
	db, err := sql.Open(dbType, dbFilename)
	if err != nil {
		return nil, err
	}

	err = createReportsTableIfNotExist(db)
	if err != nil {
		return nil, err
	}

	err = createWordReportsTableIfNotExist(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *sql.DB) error {
	return db.Close()
}

func InsertFileReport(db *sql.DB, filename string, mistakes int, notKnow int, timeSpent int64) error {
	sqlStmt := "INSERT INTO reports(date, filename, mistakes, not_know, time_spent) " +
		"VALUES(date(), '%s', %d, %d, %d);"

	_, err := db.Exec(fmt.Sprintf(sqlStmt, filename, mistakes, notKnow, timeSpent))
	if err != nil {
		return err
	}

	return nil
}

func InsertWordReport(db *sql.DB, word string, filename string, answered bool) error {
	sqlStmt := "INSERT INTO word_reports(word, filename, answered, date) " +
		"VALUES('%s', '%s', %d, date());"

	answeredColValue := 0
	if answered {
		answeredColValue = 1
	}

	_, err := db.Exec(fmt.Sprintf(sqlStmt, word, filename, answeredColValue))
	if err != nil {
		return err
	}

	return nil
}

func createReportsTableIfNotExist(db *sql.DB) error {
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

func createWordReportsTableIfNotExist(db *sql.DB) error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS word_reports (
		        word     TEXT,
				filename TEXT,
				answered INTEGER,
		        date     TEXT
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
