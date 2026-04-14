package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rahulkumarpahwa/femProject/internal/store"
	"github.com/rahulkumarpahwa/femProject/internal/tokens"
	"github.com/rahulkumarpahwa/femProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}
type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{tokenStore: tokenStore, userStore: userStore, logger: logger}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest : %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error ": "Invalid request payload"})
		return
	}

	// lets get the user
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: getUserByUsername : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "Internal Server Error"})
		return
	}

	// check for the password match
	passwordDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: passwordHash Match : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "Internal Server Error"})
		return
	}

	if !passwordDoMatch {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error ": "Invalid Credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, time.Hour*24, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: CreatingToken : %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error ": "Internal Server Error"})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
