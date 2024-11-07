package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"vsensetech.in/go_fingerprint_server/models"
	"vsensetech.in/go_fingerprint_server/payload"
)

type AuthController struct {
	AuthRepository models.AuthDetailsRepository
}

func NewAuthController(ar models.AuthDetailsRepository) *AuthController {
	return &AuthController{
		AuthRepository: ar,
	}
}

func (ac *AuthController) RegisterController(w http.ResponseWriter, r *http.Request) {
	var urlPath = strings.Split(r.URL.Path, "/")[1]
	token, err := ac.AuthRepository.Register(&r.Body, urlPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	
	// Set the cookie if the user is an admin
	if urlPath == "admin" {
		cookie := http.Cookie{
			Name:        "token",
			Value:       token,
			Expires:     time.Now().Add(24 * 365 * time.Hour), // Expires in 1 year
			Secure:      true,  // Set to true only if you're using HTTPS
			SameSite:    http.SameSiteNoneMode, // SameSite=None should be used if cross-site
			Path:        "/",
			HttpOnly:    true,
			// If needed, partition can be removed for compatibility unless specific use
			Partitioned: true, 
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
	} else {
		// Non-admin response with token data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: token})
	}
}

func (ac *AuthController) LoginController(w http.ResponseWriter, r *http.Request) {
	var urlPath = strings.Split(r.URL.Path, "/")[1]
	token, err := ac.AuthRepository.Login(&r.Body, urlPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: err.Error()})
		return
	}
	
	// Set the cookie if the user is an admin
	if urlPath == "admin" {
		cookie := http.Cookie{
			Name:        "token",
			Value:       token,
			Expires:     time.Now().Add(24 * 365 * time.Hour), // Expires in 1 year
			Secure:      true,  // Set to true only if you're using HTTPS
			SameSite:    http.SameSiteNoneMode, // SameSite=None should be used if cross-site
			Path:        "/",
			HttpOnly:    true,
			// If needed, partition can be removed for compatibility unless specific use
			Partitioned: true, 
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
		return
	} else {
		// Non-admin response with token data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: token})
	}
}

