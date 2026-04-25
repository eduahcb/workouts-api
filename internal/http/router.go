package http

import (
	"github.com/eduahcb/workouts/internal/app"
	"github.com/go-chi/chi/v5"
)

func NewRouter(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	router.Route("/workouts", func(r chi.Router) {
		r.Post("/", app.WorkoutHandler.Create)
		r.Get("/{id}", app.WorkoutHandler.GetByID)
		r.Put("/{id}", app.WorkoutHandler.Update)
		r.Delete("/{id}", app.WorkoutHandler.Delete)
	})

	return router
}
