package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	DB, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("DB : open %w\n", err)
	}

	// ping the database to check, the connection is maintained or not!
	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	// TODO : Add enhanced configuration to the connection pool settings with: db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()

	fmt.Println("Connected to Database Successfully!")
	return DB, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)

	defer func(){
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("Migrate : %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Goose Up : %w", err)
	}
	return nil
}
