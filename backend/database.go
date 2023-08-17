package backend

import (
	"database/sql"
	"embed"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

const dbDirName = "backend"
const dbFileName = "QandA.db"

//go:embed dist/massiveEdited.db
var database embed.FS

func ReadData() [][]string {
	csvFileName := filepath.Join(dbDirName, "QandA.csv") // Adjusted the path
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

func ReadDataHuge() [][]string {
	csvFileName := filepath.Join(dbDirName, "QandA.csv") // Adjusted the path
	csvInput, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("open %s failed: %s", csvFileName, err)
	}
	defer csvInput.Close()

	csvReader := csv.NewReader(csvInput)
	csvReader.FieldsPerRecord = -1

	originalRecords, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("csvReader.ReadAll failed: %s", err)
	}

	// Trim double quotes from the strings
	for i := range originalRecords {
		for j := range originalRecords[i] {
			originalRecords[i][j] = strings.Trim(originalRecords[i][j], "'")
		}
	}

	// Create a new slice to hold the multiplied data
	var allRecords [][]string

	// Append the original records 5000 times
	for i := 0; i < 5000; i++ {
		allRecords = append(allRecords, originalRecords...)
	}

	return allRecords
}

func DeleteDatabase() error {
	fn := filepath.Join(dbDirName, dbFileName)
	// Check if the database file exists
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		// Delete the existing database file
		err := os.Remove(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDatabaseIfNeeded(forceRebuild bool) error {
	fn := filepath.Join(dbDirName, dbFileName)

	if forceRebuild {
		err := DeleteDatabase()
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
		allRecords := ReadData()

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

func TermExistsOriginal(term string) (string, bool, error) {
	fn := filepath.Join(dbDirName, dbFileName)

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return "", false, fmt.Errorf("sql.Open failed: %s", err)
	}
	defer db.Close()

	var answer string
	err = db.QueryRow("SELECT answer FROM qat WHERE question = ?", term).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return answer, true, nil
}
func TermExistsEmbed(term string) (string, bool, error) {

	// Read the embedded database into memory
	dbData, err := fs.ReadFile(database, "dist/massiveEdited.db")
	if err != nil {
		return "", false, fmt.Errorf("failed to read embedded database: %s", err)
	}

	// Use the in-memory database driver of SQLite
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return "", false, fmt.Errorf("sql.Open failed: %s", err)
	}
	defer db.Close()

	// Create a temporary file in the system's temp directory
	tempFile, err := os.CreateTemp(os.TempDir(), "tempdb-*.sqlite")
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp database file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the embedded database data to the temporary file
	_, err = tempFile.Write(dbData)
	if err != nil {
		return "", false, fmt.Errorf("failed to write to temp database file: %s", err)
	}
	tempFile.Close()

	// Attach the temporary file to the in-memory database
	_, err = db.Exec("ATTACH DATABASE ? AS memdb", tempFile.Name())
	if err != nil {
		return "", false, err
	}

	var answer string
	err = db.QueryRow("SELECT answer FROM qat WHERE question = ?", term).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return answer, true, nil
}

func CheckDirectoryOld(term string) (string, bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Error getting directory: %s", err), true, nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Sprintf("Error reading files in directory: %s", err), true, nil
	}

	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	return fmt.Sprintf("Current Directory: %s\nFiles:\n%s", dir, strings.Join(fileList, "\n")), true, nil
}
func CheckDirectory(term string) (string, bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Error getting directory: %s", err), true, err
	}

	path, err := findFile(dir, "hero-image.png")
	if err != nil {
		return fmt.Sprintf("Error finding hero-image.png: %s", err), true, err
	}

	if path == "" {
		return "hero-image.png not found", true, nil
	}

	return fmt.Sprintf("Found hero-image.png at: %s", path), true, nil
}

// findFile recursively searches for a file starting from the given directory.
func findFile(dir, filename string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		// If there's a permission error, skip this directory and continue
		if os.IsPermission(err) {
			return "", nil
		}
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			// If it's a directory, search inside it
			foundPath, err := findFile(filepath.Join(dir, file.Name()), filename)
			if err != nil {
				return "", err
			}
			if foundPath != "" {
				return foundPath, nil
			}
		} else if file.Name() == filename {
			// If the file matches the filename we're looking for, return its path
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", nil
}
func ListEmbeddedFiles() (string, error) {
	return listFilesInDir(".")
}

func listFilesInDir(dir string) (string, error) {
	entries, err := database.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded directory %s: %w", dir, err)
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no entries found in directory %s", dir)
	}

	var fileList []string
	for _, entry := range entries {
		fmt.Printf("Processing entry: %s\n", entry.Name()) // Debug print
		if entry.IsDir() {
			subDirFiles, err := listFilesInDir(dir + "/" + entry.Name())
			if err != nil {
				return "", err
			}
			fileList = append(fileList, subDirFiles)
		} else {
			fileList = append(fileList, dir+"/"+entry.Name())
		}
	}

	return strings.Join(fileList, "\n"), nil
}
