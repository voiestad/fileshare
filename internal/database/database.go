package database

import (
	"database/sql"
	"mime/multipart"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
)

var db *sql.DB

func Init() error {
	os.Mkdir("database", 0775)
	var err error
	db, err = sql.Open("sqlite", "database/fileshare.db")
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS files (
            id TEXT PRIMARY KEY,
            name TEXT,
            uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	return err
}

func GetFile(id uuid.UUID) *sql.Row {
	return db.QueryRow("SELECT name FROM files WHERE id = ?", id)
}

func GetFiles() (*sql.Rows, error) {
	return db.Query("SELECT * FROM files ORDER BY name COLLATE NOCASE")
}

func AddFile(file *multipart.FileHeader) (uuid.UUID, error) {
	id := uuid.New()
	_, err := db.Exec("INSERT INTO files (id, name) VALUES (?, ?)", id, file.Filename)
	return id, err
}

func RemoveFile(id uuid.UUID) {
	db.Exec("DELETE FROM files WHERE id = ?", id)
}
