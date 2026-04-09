package store

import (
	"database/sql"
	"log"
	"time"
)

type Workout struct {
	ID int `json:"id"`
	//add the user id later here.
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	CaloriesBurned  int       `json:"calories_burned"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WorkoutEntry struct {
	ID int `json:"id"`
	// WorkoutId       int       `json:"workout_id"`
	ExerciseName    string   `json:"exercise_name"`
	Description     string   `json:"description"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds int      `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPostgresWorkoutStore(db *sql.DB, logger *log.Logger) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db, logger: logger}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(int int64) (*Workout, error)
}

func (pgws *PostgresWorkoutStore) CreateWorkout(*Workout) (*Workout, error) {

	return nil, nil
}

func (pgws *PostgresWorkoutStore) GetWorkoutById(int int64) (*Workout, error) {
	return nil, nil
}
