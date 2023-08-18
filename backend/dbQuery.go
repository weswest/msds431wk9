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

func TermExistsEmbed(term string) (string, bool, error) {

	// Read the embedded database into memory
	dbData, err := fs.ReadFile(database, filepath.Join(dbDirName, dbFileName))
	if err != nil {
		return "", false, fmt.Errorf("No database found. Did you build it? Check the following:\n1. Go to ./setup and run \"runThisFirstDatabaseCreate\". You may need to add the flag -forceRebuild=true\n2. Verify that there is a valid file at ./backend/dist/QandA.db\n3. Run 'wails build' again to rebuild.")
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
