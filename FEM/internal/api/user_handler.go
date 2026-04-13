package api

import (
	"errors"
	"log"
	"regexp"

	"github.com/rahulkumarpahwa/femProject/internal/store"
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
