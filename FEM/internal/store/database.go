package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
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
