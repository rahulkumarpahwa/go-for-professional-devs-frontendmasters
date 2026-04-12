package store

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5434 sslmode=disable")
	if err != nil {
		t.Fatalf("Opening Test DB : %v", err.Error())
	}

	// Run the migrations for our test DB
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("Migrating Test DB : %v", err.Error())
	}

	_, err = db.Exec(`TRUNCATE workouts,   workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Truncating Test DB : %v", err.Error())
	}

	return db
}
