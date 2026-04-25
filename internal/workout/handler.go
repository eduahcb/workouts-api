package workout

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eduahcb/workouts/pkg/utils"
	"github.com/jackc/pgx/v5"
)

type WorkoutHandler struct {
	workoutStore WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": err.Error()})
		return
	}

	workout, err := wh.workoutStore.GetByID(workoutId)

	if err != nil {
		wh.logger.Printf("Error fetching workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"msg": "Internal Server Error"})
		return
	}

	if workout == nil {
		wh.logger.Printf("Workout not found: ID %d", workoutId)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"msg": "workout not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": "Invalid request body"})
		return
	}

	newWorkout, err := wh.workoutStore.Create(&workout)

	if err != nil {
		wh.logger.Printf("Error creating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"msg": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": newWorkout})
}

func (wh *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	_, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": err.Error()})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		return
	}

	messages, err := workout.Valitate()
	if err != nil {
		wh.logger.Printf("Validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": messages})
		return
	}

	err = wh.workoutStore.Update(&workout)

	if err == pgx.ErrNoRows {
		wh.logger.Printf("Workout not found for update: ID %d", workout.ID)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"msg": "workout not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("Error updating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"msg": "Internal Server Error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
