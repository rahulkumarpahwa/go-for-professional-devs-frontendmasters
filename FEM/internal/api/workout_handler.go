package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rahulkumarpahwa/femProject/internal/store"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{store: store, logger: logger}
}

func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	workout, err := h.store.GetWorkoutById(workoutId)
	if err != nil {
		h.logger.Printf("Get workout Error : %v", err)
		http.Error(w, "Failed to Get the workout!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		h.logger.Printf("Create workout Error : %v", err)
		http.Error(w, "Failed to Create the workout!", http.StatusInternalServerError)
		return
	}

	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("Store workout Error : %v", err)
		http.Error(w, "Failed to Store the workout!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "created the workout\n")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (h *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := h.store.GetWorkoutById(workoutId)
	if err != nil {
		h.logger.Printf("Store workout Error : %v", err)
		http.Error(w, "Failed to fetch the workout!", http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// at this point we assume that we have found the existing workout.
	var updateWorkOutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkOutRequest)
	if err != nil {
		h.logger.Printf("Not able to decode the updated workout: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateWorkOutRequest.Title != nil {
		existingWorkout.Title = *updateWorkOutRequest.Title
	}

	if updateWorkOutRequest.Description != nil {
		existingWorkout.Description = *updateWorkOutRequest.Description
	}

	if updateWorkOutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkOutRequest.DurationMinutes
	}

	if updateWorkOutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkOutRequest.CaloriesBurned
	}

	if updateWorkOutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkOutRequest.Entries
	}

	err = h.store.UpdateWorkout(existingWorkout)
	if err != nil {
		h.logger.Printf("Not able to update workout: %v", err)
		http.Error(w, "Not able to update workout!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Workout Updated Successfully!")
}

func (h *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = h.store.DeleteWorkout(workoutId)
	if err != nil {
		h.logger.Printf("Not able to delete workout: %v", err)
		http.Error(w, "Not able to delete workout!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Workout Deleted Successfully!")

}
