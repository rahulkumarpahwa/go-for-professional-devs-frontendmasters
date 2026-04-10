package api

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/rahulkumarpahwa/femProject/internal/store"
	"github.com/rahulkumarpahwa/femProject/internal/utils"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{store: store, logger: logger}
}

func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: readIDParam : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "invalid get workout id"})
		return
	}
	workout, err := h.store.GetWorkoutById(workoutId)
	if err != nil {
		h.logger.Printf("Error: getWorkoutByID:  %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		h.logger.Printf("Error: decodeCreateWorkout : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "Unable to Decode the workout"})
		return
	}

	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("Error: createWorkout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "unable to Store workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"CreatedWorkout": createdWorkout})
}

func (h *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: readIDParam : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "invalid update workout id"})
		return
	}

	existingWorkout, err := h.store.GetWorkoutById(workoutId)
	if err != nil {
		h.logger.Printf("Error: getWorkoutById : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "unable to get the Existing Workout!"})
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
		h.logger.Printf("Error: decodeUpdatedWorkout : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "unable to decode workout"})
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
		h.logger.Printf("Error: updateWorkout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "unable to store the updated workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"Message": "Workout Updated Successfully!"})
}

func (h *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: readIDParam : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "invalid delete workout id"})
		return
	}

	err = h.store.DeleteWorkout(workoutId)
	if err != nil {
		h.logger.Printf("Error: deleteWorkout : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "Unable to Delete Workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"message": "Workout Deleted Successfully!"})
}
