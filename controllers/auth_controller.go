package controllers

import (
	"encoding/json"
	"fmt"
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
	if urlPath == "admin" {
		cookie := http.Cookie{
			Name:        "token",
			Value:       token,
			Secure:      true,
			Domain:      ".biometric.adminpanel.vsensetech.in",
			SameSite:    http.SameSiteNoneMode,
			Path:        "/",
			HttpOnly:    true,
			Expires:     time.Now().Add(24 * time.Hour),
			Partitioned: true,
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
		c := &http.Cookie{
			Name:		"token",
			Value:		"bar",
			Expires:	time.Now().Add(1 * time.Hour),
			Domain:		".godev.local",	// edit (or omit)
			Path:		"/",		// ^ ditto
			HttpOnly:	true,
		}
		fmt.Fprintln(w, "Hello world")
		http.SetCookie(w, c)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SimpleSuccessPayload{Message: "Success"})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload.SuccessPayloadWithData{Message: "Success", Data: token})
		return
	}
}
