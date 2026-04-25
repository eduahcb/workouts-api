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
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"msg": "workout not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var workoutRequest WorkoutRequest
	var workout Workout

	err := json.NewDecoder(r.Body).Decode(&workoutRequest)

	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": "Invalid request body"})
		return
	}

	workout.Title = workoutRequest.Title
	workout.Description = workoutRequest.Description
	workout.DurationMinutes = workoutRequest.DurationMinutes
	workout.CaloriesBurned = workoutRequest.CaloriesBurned
	workout.Entries = workoutRequest.Entries

	newWorkout, err := wh.workoutStore.Create(&workout)

	if err != nil {
		wh.logger.Printf("Error creating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"msg": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": newWorkout})
}

func (wh *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	var workoutRequet WorkoutRequest
	var workout Workout

	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": err.Error()})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&workoutRequet)
	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		return
	}

	messages, err := workoutRequet.Valitate()
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": messages})
		return
	}

	workout.Title = workoutRequet.Title
	workout.Description = workoutRequet.Description
	workout.DurationMinutes = workoutRequet.DurationMinutes
	workout.CaloriesBurned = workoutRequet.CaloriesBurned
	workout.Entries = workoutRequet.Entries
	workout.ID = workoutId

	err = wh.workoutStore.Update(&workout)

	if err == pgx.ErrNoRows {
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

func (wh *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paramId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"msg": err.Error()})
	}

	err = wh.workoutStore.Delete(paramId)

	if err == pgx.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"msg": "workout not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("Error deleting workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"msg": "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
