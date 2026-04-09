package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rahulkumarpahwa/femProject/internal/api"
	"github.com/rahulkumarpahwa/femProject/internal/store"
	"github.com/rahulkumarpahwa/femProject/migrations"
)

type Application struct {
	WorkoutHandler *api.WorkoutHandler
	Logger         *log.Logger
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	// Our stores will go here
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS ,".") // . means the root fof the FS and migrations is here the package and FS is embedded variable we created in fs.go
	if err !=nil{
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(pgDB, logger)

	app := &Application{Logger: logger, WorkoutHandler: workoutHandler, DB: pgDB}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available!")
}
