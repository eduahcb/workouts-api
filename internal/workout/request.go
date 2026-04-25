package workout

import (
	"github.com/eduahcb/workouts/pkg/utils"
)

type WorkoutRequest struct {
	ID              int64          `json:"id"`
	Title           string         `json:"title" validate:"required"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes" validate:"required"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

func (w *WorkoutRequest) Validate() (utils.ValidateMessageError, error) {
	return utils.ValidateStruct(w)
}
