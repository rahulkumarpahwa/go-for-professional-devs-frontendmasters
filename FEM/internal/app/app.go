package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rahulkumarpahwa/femProject/internal/api"
)

type Application struct {
	WorkoutHandler *api.WorkoutHandler
	Logger         *log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our stores will go here

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{Logger: logger, WorkoutHandler: workoutHandler}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available!")
}