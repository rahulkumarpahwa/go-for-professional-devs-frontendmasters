package store

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Workout struct {
	ID int `json:"id"`
	//add the user id later here.
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	WorkoutId       int      `json:"workout_id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
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
	UpdateWorkout(*Workout) error
	DeleteWorkout(int int64) error
}

func (pgws *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {

	tx, err := pgws.db.Begin() // begin transaction
	if err != nil {
		return nil, err
	}
	// we can't commit after rollback.
	defer tx.Rollback()

	query := `INSERT INTO workouts (title, description, duration_minutes, calories_burned) 
	VALUES ($1, $2, $3, $4) RETURNING id;`

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// we also to do entries
	for _, entry := range workout.Entries {
		query := `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

		err := tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pgws *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	tx, err := pgws.db.Begin() // begin transaction
	if err != nil {
		return nil, err
	}
	// we can't commit after rollback.
	defer tx.Rollback()

	query := `SELECT id, title, description, duration_minutes, calories_burned, created_at, updated_at FROM workouts WHERE id = $1;`

	row := tx.QueryRow(query, id)

	workout := &Workout{}
	err = row.Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// we also to do entries
	entriesQuery := `SELECT id, workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index FROM workout_entries WHERE workout_id = $1 
	ORDER BY order_index
	`
	rows, err := tx.Query(entriesQuery, id)
	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entries []WorkoutEntry
	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(&entry.ID, &entry.WorkoutId, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	workout.Entries = append(workout.Entries, entries...)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pgws *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {

	tx, err := pgws.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4 WHERE id = $5`

	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("No Rows Affected!")
	}

	// update entries by Putting the new in the place of old instead of updating each one.

	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1 `, workout.ID)
	if err != nil {
		return nil
	}

	for _, entry := range workout.Entries {
		query := `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := tx.Exec(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pgws *PostgresWorkoutStore) DeleteWorkout(id int64) error {

	tx, err := pgws.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `DELETE FROM workouts WHERE id = $1`

	result, err := tx.Exec(query, id)
	if err != nil {
		return err

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No Rows Affected!")
	}

	return tx.Commit()
}
