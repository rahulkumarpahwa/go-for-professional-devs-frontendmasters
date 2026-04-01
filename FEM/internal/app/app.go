package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rahulkumarpahwa/femProject/internal/api"
	"github.com/rahulkumarpahwa/femProject/internal/store"
)

type Application struct {
	WorkoutHandler *api.WorkoutHandler
	Logger         *log.Logger
	DB  *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our stores will go here
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(pgDB, logger)

	app := &Application{Logger: logger, WorkoutHandler: workoutHandler, DB:  pgDB}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available!")
}
