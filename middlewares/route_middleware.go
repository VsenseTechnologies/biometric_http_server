package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"vsensetech.in/go_fingerprint_server/payload"
)

func RouteMiddleware(authHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set common response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "https://biometric.adminpanel.vsensetech.in")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")

		// Get the first part of the path to determine if route needs authentication
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) > 1 && (pathParts[1] == "admin" || pathParts[1] == "users") {
			authHandler.ServeHTTP(w, r)
			return
		}

		// Return error if the route doesn't match
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: "Invalid Route"})
	})
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// JWT secret key
		jwtSecretKey := []byte("vsense")

		// Split the URL path and check for login/register routes
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) > 2 && (pathParts[2] == "login" || pathParts[2] == "register") {
			// Skip the JWT check for login and register routes
			next.ServeHTTP(w, r)
			return
		}

		// Check if the path requires authorization (e.g., admin route)
		if len(pathParts) > 1 && pathParts[1] == "admin" {
			// Retrieve JWT token from cookies
			cookie, err := r.Cookie("token")
			if err != nil {
				respondWithUnauthorized(w, "Unauthorized - no token provided")
				return
			}

			tokenString := cookie.Value

			// Parse and validate the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecretKey, nil
			})

			if err != nil || !token.Valid {
				respondWithUnauthorized(w, "Unauthorized - invalid token")
				return
			}
		}

		// If no JWT validation is needed or it's valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// Utility function to respond with an Unauthorized error in a consistent format
func respondWithUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: message})
}
