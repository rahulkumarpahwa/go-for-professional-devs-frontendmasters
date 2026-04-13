package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/rahulkumarpahwa/femProject/internal/store"
	"github.com/rahulkumarpahwa/femProject/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	store  store.UserStore
	logger *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{store: store, logger: logger}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("Username is required")
	}
	if req.Email == "" {
		return errors.New("Email is required")
	}
	// we can have the password validations as well
	if req.Password == "" {
		return errors.New("Password is required")
	}
	if req.Bio == "" {
		return errors.New("Bio is required")
	}

	if len(req.Username) > 50 {
		return errors.New("Username can't be greater than 50 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(req.Email) {
		return errors.New("Invalid Email!")
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR : decoding register request: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid register request payload!"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		h.logger.Printf("ERROR : validating register request: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		h.logger.Printf("ERROR : required Bio!: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Required Bio!"})
		return
	}

	h.store.CreateUser(user)

}
