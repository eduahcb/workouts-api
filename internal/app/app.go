package app

import (
	"log"
	"net/http"
	"os"

	"github.com/eduahcb/workouts/internal/database"
	"github.com/eduahcb/workouts/internal/workout"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *workout.WorkoutHandler
	DB             *pgxpool.Pool
}

func NewApp() (*Application, error) {
	logger := log.New(os.Stdout, "[LOG_WORKOUT]:", log.Ldate|log.Ltime)

	db, err := database.Open(logger)
	if err != nil {
		return nil, err
	}

	err = database.MigrateFS(db, database.FS, "migrations", logger)

	if err != nil {
		return nil, err
	}

	// here we can initialize all stores
	workoutStore := workout.NewPostgresWorkoutStore(db)

	// here we can initialize all handlers
	workoutHandler := workout.NewWorkoutHandler(workoutStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             db,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Status is available\n"))
}
