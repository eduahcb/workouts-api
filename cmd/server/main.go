package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eduahcb/workouts/internal/app"
	appHttp "github.com/eduahcb/workouts/internal/http"
)

func main() {
	app, err := app.NewApp()

	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	port := 8080

	app.Logger.Printf("Starting server on port: %d", port)

	router := appHttp.NewRouter(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	server.ListenAndServe()
}
