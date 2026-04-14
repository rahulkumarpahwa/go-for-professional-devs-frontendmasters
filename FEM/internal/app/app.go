package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rahulkumarpahwa/femProject/internal/api"
	"github.com/rahulkumarpahwa/femProject/internal/middleware"
	"github.com/rahulkumarpahwa/femProject/internal/store"
	"github.com/rahulkumarpahwa/femProject/migrations"
)

type Application struct {
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	Logger         *log.Logger
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	// Our stores will go here
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".") // . means the root fof the FS and migrations is here the package and FS is embedded variable we created in fs.go
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// creating the WorkoutStore
	workoutStore := store.NewPostgresWorkoutStore(pgDB, logger)
	userStore := store.NewPostgresUserStore(pgDB, logger)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userHandler := api.NewUserHandler(userStore, logger)

	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{Logger: logger, WorkoutHandler: workoutHandler, UserHandler: userHandler, TokenHandler: tokenHandler, Middleware: middlewareHandler , DB: pgDB}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available!")
}
