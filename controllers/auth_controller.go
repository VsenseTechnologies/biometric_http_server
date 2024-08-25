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
	AuthRepository models.AuthRepository
}

func NewAuthController(ar models.AuthRepository) *AuthController {
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
	if urlPath == "admin" {
		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/", // Ensure the cookie is valid site-wide
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode, // Adjust based on your use case
			Secure:   true,
			Partitioned: true,// Set to true if using HTTPS
			Expires: time.Now().Add(24 * 365 * time.Hour).UTC(),
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: token})
		return
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
	if urlPath == "admin" {
		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/", // Ensure the cookie is valid site-wide
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode, // Adjust based on your use case
			Secure:   true, // Set to true if using HTTPS
			Partitioned: true,
			Expires: time.Now().Add(24 * 365 * time.Hour).UTC(),
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: token})
		return
	}
}