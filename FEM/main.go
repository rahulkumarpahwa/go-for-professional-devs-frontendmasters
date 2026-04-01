package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/rahulkumarpahwa/femProject/internal/app"
	"github.com/rahulkumarpahwa/femProject/internal/routes"
)

func main() {
	// getting the port from the CLI flag
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse() // Important Line

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close() // close when the main end
	r := routes.SetupRoutes(app)

	Server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Logger.Printf("App Started Running on PORT %d!\n", port)

	err = Server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
