package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	DB, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("DB : open %w\n", err)
	}

	// ping the database to check, the connection is maintained or not!
	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	// Add enhanced configuration to the connection pool settings with: db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxIdleTime(5 * time.Minute)
	DB.SetConnMaxLifetime(30 * time.Minute)

	fmt.Println("Connected to Database Successfully!")
	return DB, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS) // tell goose from where to read the files

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres") // choosing the DB

	if err != nil {
		return fmt.Errorf("Migrate : %w", err)
	}

	err = goose.Up(db, dir) // what command to run (up or down) and from what folder in passed FS to read (because FS may have the multiple folder)
	if err != nil {
		return fmt.Errorf("Goose Up : %w", err)
	}
	return nil
}
