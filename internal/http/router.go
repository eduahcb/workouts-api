package http

import (
	"github.com/eduahcb/workouts/internal/app"
	"github.com/go-chi/chi/v5"
)

func NewRouter(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	router.Post("/workouts", app.WorkoutHandler.Create)
	router.Get("/workouts/{id}", app.WorkoutHandler.GetByID)
	router.Put("/workouts/{id}", app.WorkoutHandler.Update)

	return router
}
