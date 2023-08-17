package backend

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

const dbDirName = "dist"
const dbFileName = "QandA.db" //VERIFY THIS MATCHES GO EMBED BELOW

// Make sure the code below matches the dbDir and dbFile above.

//go:embed dist/QandA.db
var database embed.FS

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

func TermExistsEmbed(term string) (string, bool, error) {

	// Read the embedded database into memory
	dbData, err := fs.ReadFile(database, filepath.Join(dbDirName, dbFileName))
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
