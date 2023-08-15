package backend

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/glebarez/go-sqlite"
)

const dbDirName = "backend"
const dbFileName = "QandA.db"

func ReadData() [][]string {
	csvFileName := "QandA.csv"
	csvInput, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("open %s failed: %s", csvFileName, err)
	}
	defer csvInput.Close()

	csvReader := csv.NewReader(csvInput)
	csvReader.FieldsPerRecord = -1

	allRecords, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("csvReader.ReadAll failed: %s", err)
	}

	return allRecords
}

func CreateDatabase() error {
	// Check if the tutorDB directory already exists
	if _, err := os.Stat(dbDirName); os.IsNotExist(err) {
		allRecords := ReadData()

		fn := filepath.Join(dbDirName, dbFileName)

		db, err := sql.Open("sqlite", fn)
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(`create table if not exists qat(id integer, question text, answer text)`)
		if err != nil {
			return err
		}

		if _, err = stmt.Exec(); err != nil {
			return err
		}

		for id := 0; id < len(allRecords); id++ {
			record := allRecords[id]
			question := record[0]
			answer := record[1]

			stmt, err = db.Prepare("insert into qat(id, question, answer) values(?, ?, ?)")
			if err != nil {
				return err
			}

			_, err = stmt.Exec(strconv.Itoa(id), question, answer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ListDatabase() ([]string, error) {
	var results []string

	fn := filepath.Join(dbDirName, dbFileName)

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select * from qat")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var question string
		var answer string
		rows.Scan(&id, &question, &answer)
		results = append(results, fmt.Sprintf("The answer to %s is %s", question, answer))
	}
	return results, nil
}

func TermExists(term string) (string, bool) {
	fn := filepath.Join(dbDirName, dbFileName)

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		log.Fatalf("sql.Open failed: %s", err)
	}
	defer db.Close()

	var answer string
	err = db.QueryRow("SELECT answer FROM qat WHERE question = ?", term).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		panic(err)
	}
	return answer, true
}
