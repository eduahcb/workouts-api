package workout

import (
	"github.com/eduahcb/workouts/pkg/utils"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title" validate:"required"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes" validate:"required"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

func (w *Workout) Valitate() (utils.ValidateMessageError, error) {
	return utils.ValidateStruct(w)
}
