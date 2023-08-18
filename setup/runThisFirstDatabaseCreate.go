package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

const csvDirName = "../backend"
const dbDirName = "../backend/dist"
const dbFileName = "QandA.db"

func readData() [][]string {
	csvFileName := filepath.Join(csvDirName, "QandA.csv") // Adjusted the path
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

	// Trim double quotes from the strings
	for i := range allRecords {
		for j := range allRecords[i] {
			allRecords[i][j] = strings.Trim(allRecords[i][j], "'")
		}
	}

	return allRecords
}

func deleteDatabase() error {
	fn := filepath.Join(dbDirName, dbFileName)
	// Check if the database file exists
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		// Delete the entire directory
		err := os.RemoveAll(dbDirName)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDatabaseIfNeeded(forceRebuild bool) error {
	fn := filepath.Join(dbDirName, dbFileName)

	if forceRebuild {
		err := deleteDatabase()
		if err != nil {
			return err
		}
		// Create the directory again
		err = os.MkdirAll(dbDirName, os.ModePerm)
		if err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the table 'qat' exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='qat'").Scan(&tableName)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if tableName != "qat" {
		allRecords := readData()

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
	} else {
		fmt.Println("Database table 'qat' already exists.")
	}
	return nil
}

func main() {
	// Handle forceRebuild input
	forceRebuildPtr := flag.Bool("forceRebuild", true, "Force database rebuild")
	flag.Parse()

	if *forceRebuildPtr {
		err := createDatabaseIfNeeded(true)
		if err != nil {
			log.Fatalf("Failed to create database with force rebuild: %s", err)
		}
	} else {
		fn := filepath.Join(dbDirName, dbFileName)
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			err := createDatabaseIfNeeded(false)
			if err != nil {
				log.Fatalf("Failed to create database: %s", err)
			}
		} else {
			fmt.Println("Database exists. If you want to force a rebuild, enter command with -forceRebuild=true")
		}
	}
}
