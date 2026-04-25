package workout

import (
	"github.com/eduahcb/workouts/pkg/utils"
)

type WorkoutRequest struct {
	ID              int            `json:"id"`
	Title           string         `json:"title" validate:"required"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes" validate:"required"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

func (w *WorkoutRequest) Valitate() (utils.ValidateMessageError, error) {
	return utils.ValidateStruct(w)
}
