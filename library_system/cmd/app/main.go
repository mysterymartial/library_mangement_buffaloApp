package main

import (
	"library-system/actions"
	"log"

	"github.com/gobuffalo/envy"
)

func main() {
	env := envy.Get("GO_ENV", "development")
	app := actions.App()

	if app == nil {
		log.Fatal("Failed to initialize application")
	}

	log.Printf("Starting %s server on :3000", env)

	if err := app.Serve(); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
