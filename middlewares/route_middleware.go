package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"vsensetech.in/go_fingerprint_server/payload"
)

func RouteMiddleware(authHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Set to your frontend's origin
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
