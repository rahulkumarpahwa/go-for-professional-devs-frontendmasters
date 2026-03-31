package main

import "github.com/rahulkumarpahwa/femProject/internal/app"

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("App Started Running!")
}
