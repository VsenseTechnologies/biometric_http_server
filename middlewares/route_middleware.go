package middlewares

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"

	// "github.com/golang-jwt/jwt/v5"
	"vsensetech.in/go_fingerprint_server/payload"
)

func RouteMiddleware(authHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "https://biometric-dashboard-admin.vercel.app") // Set to your frontend's origin
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		var url = strings.Split(r.URL.Path, "/")[1]
		if url == "admin" || url == "users" {
			authHandler.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(payload.SimpleFailedPayload{ErrorMessage: "Invalid Route"})
	})
}

// func JwtMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var jwtSecretKey = []byte("vsense")
// 		fmt.Println(r.URL.Path)
// 		var path string = strings.Split(r.URL.Path, "/")[2]
// 		fmt.Println(path)
// 		// Bypass the middleware for login and register routes
// 		if path == "login" || path == "register" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// Get JWT token from cookies
// 		cookie, err := r.Cookie("token")
// 		if err != nil {
// 			http.Error(w, "Unauthorized - no token provided", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := cookie.Value
// 		fmt.Println(tokenString)

// 		// Parse and validate the token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Check the signing method
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return jwtSecretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		// If token is valid, proceed to the next handler
// 		next.ServeHTTP(w, r)
// 	})
// }
